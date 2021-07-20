package db

import (
	"fmt"
	"os"

	"github.com/miranaky/kaengkaengcoin/utils"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

func getDBName() string {
	port := os.Args[2][7:]
	return fmt.Sprintf("%s_%s.db", dbName, port)

}

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDBName(), 0600, nil)
		utils.HandleErr(err)
		db = dbPointer
		err = db.Update(func(t *bolt.Tx) error {
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockChain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func EmptyBlocks() {
	DB().Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
