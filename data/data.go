// Persist data with boltdb.
// https://github.com/boltdb/bolt/blob/master/README.md

package data

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/zhengxiaoyao0716/zmodule/config"
	"github.com/zhengxiaoyao0716/zmodule/event"
)

type bName []byte

// Do .
func (n *bName) Do(fn func(*bolt.Bucket) error) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error { return fn(tx.Bucket(*n)) }
}

var (
	bMain    bName
	bAccount bName
)

var (
	db     *bolt.DB
	dbPath string
)

// Init used to create or open thr database.
func Init() {
	var err error

	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// close db before stop.
	event.On(event.KeyStop, func(event.Event) error { return db.Close() })

	err = db.Update(func(tx *bolt.Tx) error {
		initBucket := func(name string) (bName, error) {
			n := bName(name)
			_, err := tx.CreateBucketIfNotExists(n)
			if err != nil {
				return nil, fmt.Errorf("create bucket failed, name: %s, %s", name, err)
			}
			return n, nil
		}
		if bMain, err = initBucket("main"); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if bAccount, err = initBucket("account"); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	event.OnInit(func(e event.Event) error {
		event.On(event.KeyStart, func(event.Event) error {
			dbPath = config.GetString("db_path")
			Init()
			return nil
		})
		return nil
	})
}
