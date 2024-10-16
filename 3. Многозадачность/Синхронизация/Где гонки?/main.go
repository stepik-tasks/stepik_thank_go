package main

import (
	"fmt"
	"time"
)

func delay(duration time.Duration, fn func()) func() {
	canceled := false

	go func() {
		time.Sleep(duration)
		if !canceled {
			fn()
		}
	}()

	cancel := func() {
		canceled = true
	}
	return cancel
}

func main() {
	work := func() {
		fmt.Println("work done")
	}

	cancel := delay(50*time.Millisecond, work)
	time.Sleep(50 * time.Millisecond)
	go cancel()
}
