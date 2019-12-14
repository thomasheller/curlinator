package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Crawler struct {
	db *DB
}

func (c *Crawler) Run(ticker *time.Ticker, done chan bool) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			c.next()
		}
	}
}

func (c *Crawler) next() {
	url := c.db.NextURL()

	if url == "" {
		return
	}

	log.Printf("Downloading: %s", url)

	err := c.crawl(url)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	c.db.DeleteURL(url)
}

func (c *Crawler) crawl(u string) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", u, nil)

	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "")

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	filename := url.QueryEscape(u)

	out, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	return nil
}
