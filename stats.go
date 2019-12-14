package main

import (
	"fmt"
	"time"
)

type Stats struct{}

func (s *Stats) Format(count int) string {
	if count == 0 {
		return "Database is empty"
	}

	duration := time.Duration(count) * delay
	eta := time.Now().Add(duration)

	return fmt.Sprintf("%d URL(s) in database, approximately %s left (ETA: %s)", count, duration, eta.Format(time.RFC1123Z))
}
