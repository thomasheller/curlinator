package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/url"

	bolt "go.etcd.io/bbolt"
)

const dbFile = "curlinator.db"

type DB struct {
	db *bolt.DB
}

func OpenDB() (*DB, error) {
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) AddURL(u string) error {
	if u == "" {
		return errors.New("Empty URL")
	}

	parsed, err := url.Parse(u)

	if err != nil {
		return fmt.Errorf("Failed to parse URL: %v", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("Unsupported URL scheme: %s", parsed.Scheme)
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(u), []byte(""))

		if err != nil {
			log.Printf("Error adding URL: %v", err)
		}

		if err == nil {
			log.Printf("Added URL: %s", u)
		}

		return err
	})
}

func (d *DB) itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (d *DB) DeleteURL(url string) error {
	if url == "" {
		return errors.New("Empty URL")
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete([]byte(url))

		if err != nil {
			log.Printf("Error removing URL: %v", err)
		}

		if err == nil {
			log.Printf("Removed URL: %s", url)
		}

		return err
	})
}

func (d *DB) ListURLs() []string {
	var urls []string

	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		b.ForEach(func(k, v []byte) error {
			url := make([]byte, len(k))
			copy(url, k)

			urls = append(urls, string(url))

			// fmt.Printf("key=%s, value=%s\n", k, v)
			fmt.Printf("%s\n", k)
			return nil
		})

		return nil
	})

	return urls
}

func (d *DB) NextURL() string {
	var url []byte

	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		c := b.Cursor()

		k, _ := c.First()

		if k != nil {
			url = make([]byte, len(k))
			copy(url, k)
		}

		return nil
	})

	return string(url)
}

func (d *DB) Size() int {
	var count int

	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		s := b.Stats()
		count = s.KeyN
		return nil
	})

	return count
}

func (d *DB) Close() error {
	return d.db.Close()
}
