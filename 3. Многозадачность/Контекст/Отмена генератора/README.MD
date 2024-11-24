**Отмена генератора**

Есть функция `generate()`, которая генерит числа:

```go
// генерит целые числа от start и до бесконечности
func generate(cancel <-chan struct{}, start int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := start; ; i++ {
            select {
            case out <- i:
            case <-cancel:
                return
            }
        }
    }()
    return out
}
```

Функция использует канал отмены. Переделайте на контекст:

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    generated := generate(ctx, 11)
    for num := range generated {
        fmt.Print(num, " ")
        if num > 14 {
            break
        }
    }
    fmt.Println()
}
```

```text
11 12 13 14 15
```