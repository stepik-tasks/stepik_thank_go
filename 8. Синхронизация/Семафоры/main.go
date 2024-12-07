// Пишем семафор
package main

import (
	"fmt"
	"sync"
	"time"
)

// начало решения

// Semaphore представляет семафор синхронизации.
type Semaphore struct {
	pool chan int
}

// NewSemaphore создает новый семафор указанной вместимости.
func NewSemaphore(n int) *Semaphore {
	ch := make(chan int, n)
	i := 0
	for i < n {
		ch <- 1
		i++
	}

	return &Semaphore{pool: ch}
}

// Acquire занимает место в семафоре, если есть свободное.
// В противном случае блокирует вызывающую горутину.
func (s *Semaphore) Acquire() {
	<-s.pool
}

// TryAcquire занимает место в семафоре, если есть свободное,
// и возвращает true. В противном случае просто возвращает false.
func (s *Semaphore) TryAcquire() bool {
	select {
	case <-s.pool:
		return true
	default:
		return false
	}
}

// Release освобождает место в семафоре и разблокирует
// одну из заблокированных горутин (если такие были).
func (s *Semaphore) Release() {
	s.pool <- 1
}

// конец решения

func main() {
	const maxConc = 4
	sema := NewSemaphore(maxConc)
	start := time.Now()

	const nCalls = 12
	var wg sync.WaitGroup
	wg.Add(nCalls)

	for i := 0; i < nCalls; i++ {
		sema.Acquire()
		go func() {
			defer wg.Done()
			defer sema.Release()
			time.Sleep(10 * time.Millisecond)
			fmt.Print(".")
		}()
	}

	wg.Wait()

	fmt.Printf("\n%d calls took %d ms\n", nCalls, time.Since(start).Milliseconds())
	/*
		............
		12 calls took 30 ms
	*/
}
