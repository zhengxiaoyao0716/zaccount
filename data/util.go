package data

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/zhengxiaoyao0716/zaccount/helper/encrypt"
	"github.com/zhengxiaoyao0716/zmodule/info"
)

var resetPepperQueue = []func(tx *bolt.Tx, r func([]byte) ([]byte, error)) error{}

// ResetPepper .
func ResetPepper(old, new string) error {
	pepper := func(key string) []byte {
		if key == "" {
			return encrypt.DefaultPepper
		}
		return encrypt.Pepper(key)
	}
	oldPepper := pepper(old)
	newPepper := pepper(new)
	return db.Update(func(tx *bolt.Tx) error {
		// Check and reset the papper
		if err := bMain.Do(func(b *bolt.Bucket) error {
			var (
				key   = []byte("pepper_check")
				plain = []byte("pepper_check")
			)

			secret := b.Get(key)
			if secret != nil {
				restore, err := encrypt.AesDecrypt(secret, oldPepper)
				if err != nil {
					return err
				}
				if bytes.Compare(plain, restore) != 0 {
					return errors.New("reset pepper failed, old pepper not match")
				}
			}

			secret, err := encrypt.AesEncrry(plain, newPepper)
			if err != nil {
				return err
			}
			if err := b.Put(key, secret); err != nil {
				return err
			}
			return nil
		})(tx); err != nil {
			return err
		}
		var errStr string
		// Reset all data encrypt with old pepper
		for _, reset := range resetPepperQueue {
			if err := reset(tx, func(old []byte) ([]byte, error) {
				restore, err := encrypt.AesDecrypt(old, oldPepper)
				if err != nil {
					errStr += err.Error()
					return nil, err
				}
				new, err := encrypt.AesEncrry(restore, newPepper)
				if err != nil {
					errStr += err.Error()
					return nil, err
				}
				return new, nil
			}); err != nil {
				errStr += err.Error()
			}
		}
		return nil
	})
}

// BackupHandleFunc is a HTTP handle used to backup database over HTTP
// https://github.com/boltdb/bolt#database-backups
func BackupHandleFunc(w http.ResponseWriter, req *http.Request) {
	err := db.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.db"`, info.Name()))
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
