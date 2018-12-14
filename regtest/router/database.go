package router

import (
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Database struct {
	DB *bolt.DB
}

// New creates or opens a db
func (bb *Database) New() error {
	db, err := bolt.Open("Utxos.db", 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}

	bb.DB = db

	return nil
}

// Update writes
func (bb *Database) Update(addrBucket string, txBucket string, utxoBucket string, key string, value string) error {
	log.Println("********************")
	log.Println("ADDING NEW ENTRY IN DATABASE:")
	log.Println("Address bucket:", addrBucket)
	log.Println("Tx bucket:", txBucket)
	log.Println("Utxo bucket:", utxoBucket)
	log.Println("Key:", key)
	log.Println("Value:", value)

	return bb.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(addrBucket))
		if err != nil {
			return err
		}
		sb, err := b.CreateBucketIfNotExists([]byte(txBucket))
		if err != nil {
			return err
		}
		ssb, err := sb.CreateBucketIfNotExists([]byte(utxoBucket))
		if err != nil {
			return err
		}

		err = ssb.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})
}

// Get retrieves an entry from the bucket
func (bb *Database) Get(addrBucket string, txBucket string, utxoBucket string, key string) (string, error) {
	var response string
	err := bb.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(addrBucket))
		if b == nil {
			return fmt.Errorf("Bucket %s not found", addrBucket)
		}
		sb := b.Bucket([]byte(txBucket))
		if sb == nil {
			return fmt.Errorf("Bucket %s not found", txBucket)
		}
		ssb := sb.Bucket([]byte(utxoBucket))
		if ssb == nil {
			return fmt.Errorf("Bucket %s not found", utxoBucket)
		}
		response = string(ssb.Get([]byte(key)))

		return nil
	})

	return response, err
}

// List retrieves all data in a bucket
func (bb *Database) List(addrBucket string, txBucket string) ([][]byte, error) {
	var response [][]byte
	err := bb.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(addrBucket))
		if b == nil {
			return fmt.Errorf("Bucket %s not found", addrBucket)
		}
		sb := b.Bucket([]byte(txBucket))
		if sb == nil {
			return fmt.Errorf("Bucket %s not found", txBucket)
		}

		c := sb.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			ssb := sb.Bucket(k)
			ssb.ForEach(func(_, value []byte) error {
				response = append(response, value)
				return nil
			})
		}

		return nil
	})

	return response, err
}

// Delete an entry from the bucket
func (bb *Database) Delete(addrBucket string, txBucket string, utxoBucket string, key string) error {
	return bb.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(addrBucket))
		if b == nil {
			return fmt.Errorf("Bucket %s not found", addrBucket)
		}
		sb := b.Bucket([]byte(txBucket))
		if sb == nil {
			return fmt.Errorf("Bucket %s not found", txBucket)
		}
		ssb := sb.Bucket([]byte(utxoBucket))
		if ssb == nil {
			return fmt.Errorf("Bucket %s not found", utxoBucket)
		}
		return ssb.Delete([]byte(key))
	})
}

// Close ends connection to the db
func (bb *Database) Close() error {
	return bb.DB.Close()
}
