package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/chetan-99/account-cred-manager-go-grpc/internal/config"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/utils"
	"github.com/dgraph-io/badger"
)

type DbHandler struct {
	*badger.DB
}

func NewBadgerDB(cfg *config.AppConfig) *DbHandler {
	bdb, err := badger.Open(badger.DefaultOptions(cfg.DB_PATH).WithLogger(nil))
	if err != nil {
		log.Fatalf("Failed to initialize badger db - err - %v", err)
	}

	db_handler := &DbHandler{
		DB: bdb,
	}

	return db_handler
}

func (db *DbHandler) Add_KV(account_id int32, data []byte) error {
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(utils.Convert_int32_to_byte(account_id), data)
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully added key account id - %v\n", account_id)
	return nil
}

func (db *DbHandler) Get(account_id int32) ([]byte, error) {
	var value []byte

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(utils.Convert_int32_to_byte(account_id))
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		return []byte{}, err
	}
	return value, nil
}

func (db *DbHandler) GetAllKeys() ([]int32, error) {
	var keys []int32

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			var shortInt32 int32
			shortReader := bytes.NewReader(k)
			err := binary.Read(shortReader, binary.BigEndian, &shortInt32)
			if err != nil {
				return err
			}
			keys = append(keys, shortInt32)
		}
		return nil
	})

	return keys, err
}
