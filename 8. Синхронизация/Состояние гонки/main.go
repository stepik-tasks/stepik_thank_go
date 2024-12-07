// Карта с синхронизацией (sync.Map)
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	counter := sync.Map{}

	go func() {
		defer wg.Done()
		for range 100 {
			count, _ := counter.LoadOrStore("hello", 0)
			counter.Store("hello", count.(int)+1)
		}
	}()

	go func() {
		defer wg.Done()
		for range 100 {
			count, _ := counter.LoadOrStore("hello", 0)
			counter.Store("hello", count.(int)+1)
		}
	}()

	wg.Wait()
	count, _ := counter.Load("hello")
	fmt.Println("hello =", count.(int))
}
