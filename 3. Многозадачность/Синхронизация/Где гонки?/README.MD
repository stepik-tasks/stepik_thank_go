🤔 **Где гонки? (флаг)**

Я написал функцию `delay()`, которая запускает переданную функцию спустя указанный интервал времени. Она возвращает функцию `cancel()`, чтобы можно было отменить запланированное выполнение. Отмена реализована через установку флага `canceled`:

```go
func delay(duration time.Duration, fn func()) func() {
    canceled := false          // (1)

    go func() {
        time.Sleep(duration)
        if !canceled {         // (2)
            fn()
        }
    }()

    cancel := func() {
        canceled = true        // (3)
    }
    return cancel              // (4) 
}

```

Пример вызова:

```go
func main() {
	work := func() {
		fmt.Println("work done")
	}

	cancel := delay(50*time.Millisecond, work)
	time.Sleep(50 * time.Millisecond)
	go cancel()
}
```

Как считаете, есть ли тут гонки? Если да, в каких строках?

