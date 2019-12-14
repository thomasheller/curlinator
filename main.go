package main

import (
	"log"
	"sync"
	"time"
)

const (
	bucket = "URLs"
	delay  = 2 * time.Second // TODO: make configurable
)

func main() {
	db, err := OpenDB()

	if err != nil {
		log.Fatalf("Failed to open database file: %v", err)
	}

	defer db.Close()

	log.Println("Starting")

	stats := &Stats{}

	log.Println(stats.Format(db.Size()))

	var wg sync.WaitGroup

	ticker := time.NewTicker(delay)
	done := make(chan bool)

	wg.Add(1)

	go func() {
		crawler := &Crawler{db: db}
		crawler.Run(ticker, done)
		wg.Done()
	}()

	wg.Add(1)

	go func() {
		api := &API{db: db, stats: stats}
		api.Run()
		wg.Done()
	}()

	wg.Wait()
}
