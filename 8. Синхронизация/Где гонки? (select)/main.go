// Где гонки? (select)
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// delay выполняет fn после задержки duration.
// Возвращает функцию, которая отменяет выполнение fn.
func delay(duration time.Duration, fn func()) func() {
	canceled := make(chan struct{}) // (1)

	go func() {
		timer := time.NewTimer(duration)
		select {
		case <-timer.C:
			fn() // (2)
		case <-canceled:
			timer.Stop() // (3)
		}
	}()

	return func() {
		select {
		case <-canceled:
		default:
			close(canceled) // (4)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// запускаем работу через 50 мс
	work := func() {
		fmt.Println("work done")
		wg.Done()
	}
	cancel := delay(50*time.Millisecond, work)
	defer cancel()

	// отменяем работу через 20 мс c вероятностью 50%
	time.Sleep(20 * time.Millisecond)
	if rand.Intn(2) == 0 {
		cancel()
		fmt.Println("canceled")
		wg.Done()
	}

	wg.Wait()
}
