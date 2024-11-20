package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// начало решения

type Total struct {
	data atomic.Int64
}

func (t *Total) Increment() {
	t.data.Add(1)
}

func (t *Total) Value() int {
	return int(t.data.Load())
}

// конец решения

func main() {
	var wg sync.WaitGroup

	var total Total

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				total.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Println("total", total.Value())
}
