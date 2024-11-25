**Работяга с ожиданием**

Я реализовал тип-работягу `Worker`. Он выполняет заданную функцию в цикле, пока та не сломается:

```go
// Worker выполняет заданную функцию в цикле, пока не будет остановлен.
type Worker struct {
    fn func() error
}

// NewWorker создает новый экземпляр Worker с заданной функцией.
// Но пока не запускает цикл с функцией.
func NewWorker(fn func() error) *Worker {
    return &Worker{fn: fn}
}

// Start запускает отдельную горутину, в которой циклически
// выполняет заданную функцию, пока та не вернет ошибку.
func (w *Worker) Start() {
    go func() {
        for {
            err := w.fn()
            if err != nil {
                return
            }
        }
    }()
}
```

Пример использования:

```go
// Завершение по ошибке
count := 3
fn := func() error {
    fmt.Print(count, " ")
    count--
    if count == 0 {
        return errors.New("count is zero")
    }
    time.Sleep(10 * time.Millisecond)
    return nil
}

worker := NewWorker(fn)
worker.Start()
time.Sleep(25 * time.Millisecond)

// 3 2 1
```

Теперь пришло время доработать работягу под новые требования.

➊ Добавить метод Stop, который останавливает работягу:

```go
// Завершение по Stop
count := 3
fn := func() error {
    fmt.Print(count, " ")
    count--
    time.Sleep(10 * time.Millisecond)
    return nil
}

worker := NewWorker(fn)
worker.Start()
time.Sleep(25 * time.Millisecond)
worker.Stop()

// 3 2 1
```

➋ Добавить метод Wait, который блокирует клиента, пока работяга не завершится:

```go
// Ожидание завершения через Wait
count := 3
fn := func() error {
    fmt.Print(count, " ")
    count--
    time.Sleep(10 * time.Millisecond)
    return nil
}

worker := NewWorker(fn)
worker.Start()

// эта горутина остановит работягу через 25 мс
go func() {
    time.Sleep(25 * time.Millisecond)
    worker.Stop()
}()

// подождем, пока кто-нибудь остановит работягу
worker.Wait()
fmt.Println("done")

// 3 2 1 done
```

Доработайте Worker под новые требования. Не используйте каналы и контекст. Используйте группу ожидания. Можете использовать логическую переменную для статуса работяги (запущен/остановлен) — не опасаясь гонок, раз мы их еще не проходили. Детальные требования — в комментариях по коду.