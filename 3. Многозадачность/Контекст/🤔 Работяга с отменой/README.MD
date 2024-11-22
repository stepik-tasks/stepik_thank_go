🤔 **Работяга с отменой**

Я реализовал тип-работягу Worker. Он выполняет заданную функцию в цикле, пока та не сломается:

```go
// Worker выполняет заданную функцию в цикле, пока не будет остановлен.
// Гарантируется, что Worker используется только в одной горутине.
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

time.Sleep(35 * time.Millisecond)
// 3 2 1
```

Теперь пришло время доработать работягу под новые требования.

➊ Добавить метод Stop, который останавливает работягу:

```go
count := 9
fn := func() error {
    fmt.Print(count, " ")
    count--
    time.Sleep(10 * time.Millisecond)
    return nil
}

worker := NewWorker(fn)
worker.Start()
time.Sleep(105 * time.Millisecond)
worker.Stop()

fmt.Println()
// 9 8 7 6 5 4 3 2 1 0
```

➋ Добавить метод AfterStop. Он регистрирует функцию, которая будет вызвана после остановки работяги:

```go
fn := func() error { return nil }

worker := NewWorker(fn)
worker.AfterStop(func() {
    fmt.Println("called after stop")
})

worker.Start()
worker.Stop()

time.Sleep(10 * time.Millisecond)
// called after stop
```

➌ Добавить метод Err, который возвращает причину остановки (либо функция сломалась, либо остановка через Stop).

Доработайте Worker под новые требования. Подробности и нюансы — в комментариях по коду.

