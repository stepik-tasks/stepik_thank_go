// Ограничитель вызовов
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrBusy = errors.New("busy")
var ErrCanceled = errors.New("canceled")

// начало решения

// throttle следит, чтобы функция fn выполнялась не более limit раз в секунду.
// Возвращает функции handle (выполняет fn с учетом лимита) и cancel (останавливает ограничитель).
func throttle(limit int, fn func()) (handle func() error, cancel func()) {
	var counter int
	var locker = sync.RWMutex{}

	done := make(chan bool)

	go func() {
		t := time.NewTicker(time.Second * 1)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				locker.Lock()
				counter = 0
				locker.Unlock()
			case <-done:
				return
			}
		}
	}()

	cancel = func() {
		select {
		case <-done:
			return
		default:
			locker.Lock()
			counter = 0
			locker.Unlock()
			close(done)
		}
	}

	handle = func() error {
		select {
		case <-done:
			return ErrCanceled
		default:
			locker.Lock()
			available := counter < limit
			counter++
			locker.Unlock()

			if available {
				fn()
				return nil
			} else {
				return ErrBusy
			}
		}
	}

	return
}

// конец решения

func main() {
	work := func() {
		fmt.Print(".")
		time.Sleep(time.Millisecond * 500)
	}

	handle, cancel := throttle(5, work)
	defer cancel()

	const n = 200
	var nOK, nErr int
	for i := 0; i < n; i++ {
		err := handle()
		if err == nil {
			nOK += 1
		} else {
			nErr += 1
		}
	}
	fmt.Println()
	fmt.Printf("%d calls: %d OK, %d busy\n", n, nOK, nErr)
}
