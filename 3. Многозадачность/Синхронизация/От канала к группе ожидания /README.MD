Я написал функцию `timeit`. Она принимает на входе другую функцию, и выполняет ее много раз параллельно несколькими
горутинами, после чего возвращает время выполнения:

```go
// timeit выполняет nIter вызовов функции fn
// с помощью nWorkers параллельных обработчиков,
// и возвращает время выполнения в миллисекундах.
func timeit(nIter int, nWorkers int, fn func()) int {
    done := make(chan struct{}, nWorkers)
    start := time.Now()

    // работают nWorkers параллельных обработчиков
    for i := 0; i < nWorkers; i++ {
        go func() {
            // каждый обработчик выполняет nIter/nWorkers итераций
            for i := 0; i < nIter/nWorkers; i++ {
                fn()
            }
            done <- struct{}{}
        }()
    }

    // дожидаемся завершения обработчиков
    for i := 0; i < nWorkers; i++ {
        <-done
    }

    return int(time.Since(start).Milliseconds())
}
```

Пример использования:

```go
rnd := rand.New(rand.NewSource(42))

fn := func() {
    // "работа" занимает от 10 до 50 мс
    n := 10 + rnd.Intn(40)
    time.Sleep(time.Duration(n) * time.Millisecond)
}

elapsed := timeit(96, 4, fn)
fmt.Printf("96 iterations, 4 workers, took %dms\n", elapsed)
```

```go
96 iterations, 4 workers, took 896ms
```

Сейчас `timeit` использует канал завершения. Перепишите на группу ожидания.

Гарантии (чтобы не усложнять задачу):

- `nIter > 0`, `nWorkers > 0`
- `nWorkers <= nIter`
- `nIter % nWorkers == 0`
