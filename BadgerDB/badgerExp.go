package main

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

func main() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/badger"
	opts.ValueDir = "/tmp/badger"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_ = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("key"), []byte("Hello"))
		return err
	})

	var valCopy []byte
	_ = db.View(func(txn *badger.Txn) error {
		item, _ := txn.Get([]byte("key"))
		_ = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return nil
	})
	fmt.Printf("The value is: %s\n", valCopy)
}
