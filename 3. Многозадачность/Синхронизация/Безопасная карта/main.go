package main

import (
	"fmt"
	"sync"
)

// начало решения

type Counter struct {
	lock sync.Mutex
	data map[string]int
}

func (c *Counter) Increment(str string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.data[str]; ok {
		c.data[str]++
	} else {
		c.data[str] = 1
	}
}

func (c *Counter) Value(str string) int {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.data[str]; ok {
		return c.data[str]
	} else {
		return 0
	}
}

func (c *Counter) Range(fn func(key string, val int)) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for k, v := range c.data {
		fn(k, v)
	}
}

func NewCounter() *Counter {
	return &Counter{data: map[string]int{}}
}

// конец решения

func main() {
	counter := NewCounter()

	var wg sync.WaitGroup
	wg.Add(3)

	increment := func(key string, val int) {
		defer wg.Done()
		for ; val > 0; val-- {
			counter.Increment(key)
		}
	}

	go increment("one", 100)
	go increment("two", 200)
	go increment("three", 300)

	wg.Wait()

	fmt.Println("two:", counter.Value("two"))

	fmt.Print("{ ")
	counter.Range(func(key string, val int) {
		fmt.Printf("%s:%d ", key, val)
	})
	fmt.Println("}")
}
