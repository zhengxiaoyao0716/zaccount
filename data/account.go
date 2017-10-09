package data

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/zhengxiaoyao0716/zaccount/helper/encrypt"
)

// PutAccount .
func PutAccount(id, password string) func(*bolt.Tx) error {
	return bAccount.Do(func(b *bolt.Bucket) error {
		b, err := b.CreateBucketIfNotExists([]byte(id))
		if err != nil {
			return err
		}
		secret, err := encrypt.Encrypt(password)
		if err != nil {
			return err
		}
		if err := b.Put([]byte("secret"), secret); err != nil {
			return err
		}
		return nil
	})
}

// CheckAccount .
func CheckAccount(id, password string) func(*bolt.Tx) error {
	return bAccount.Do(func(b *bolt.Bucket) error {
		b = b.Bucket([]byte(id))
		if b == nil {
			return errors.New("account not found")
		}
		if err := encrypt.Compare(b.Get([]byte("secret")), password); err != nil {
			return errors.New("password not match, " + err.Error())
		}
		return nil
	})
}

// DeleteAccount .
func DeleteAccount(id string) func(*bolt.Tx) error {
	return bAccount.Do(func(b *bolt.Bucket) error { return b.DeleteBucket([]byte(id)) })
}

// Account .
func Account(id string, fn func(*bolt.Bucket) error) func(*bolt.Tx) error {
	return bAccount.Do(func(b *bolt.Bucket) error { return fn(b.Bucket([]byte(id))) })
}

func init() {
	resetPepperQueue = append(resetPepperQueue, func(tx *bolt.Tx, r func([]byte) ([]byte, error)) error {
		return bAccount.Do(func(b *bolt.Bucket) error {
			return b.ForEach(func(k, _ []byte) error {
				v := b.Bucket(k)
				old := v.Get([]byte("secret"))
				new, err := r(old)
				if err != nil {
					return err
				}
				return v.Put(k, new)
			})
		})(tx)
	})
}
