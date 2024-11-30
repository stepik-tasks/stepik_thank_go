### 🤔 Безопасная карта

Я реализовал тип `ConcMap` — безопасную карту для многозадачной среды:

```go
// ConcMap - безопасная в многозадачной среде карта.
type ConcMap[K comparable, V any] struct {
    items map[K]V
    lock  sync.Mutex
}

// NewConcMap создает новую карту.
func NewConcMap[K comparable, V any]() *ConcMap[K, V] {
    return &ConcMap[K, V]{items: map[K]V{}}
}
```

У карты есть две базовые операции — `Get` и `Set`:

```go
// Get возвращает значение по ключу.
func (cm *ConcMap[K, V]) Get(key K) V {
    cm.lock.Lock()
    defer cm.lock.Unlock()
    return cm.items[key]
}

// Set устанавливает значение по ключу.
func (cm *ConcMap[K, V]) Set(key K, val V) {
    cm.lock.Lock()
    defer cm.lock.Unlock()
    cm.items[key] = val
}
```

Пришло время добавить еще парочку.

➊ `SetIfAbsent` устанавливает новое значение по ключу и возвращает его, 
но только если такого ключа нет в карте. Если ключ уже есть — возвращает 
старое значение по ключу:

```go
var wg sync.WaitGroup
wg.Add(2)

m := NewConcMap[string, int]()

go func() {
    defer wg.Done()
    time.Sleep(5 * time.Millisecond)
    m.SetIfAbsent("hello", 42)
}()

go func() {
    defer wg.Done()
    time.Sleep(10 * time.Millisecond)
    m.SetIfAbsent("hello", 84)
}()

wg.Wait()
fmt.Println("hello =", m.Get("hello"))
// hello = 42 (от первой горутины)
```

➋ `Compute` устанавливает значение по ключу, применяя к нему функцию. Функция выполняется атомарно:

```go
var wg sync.WaitGroup
wg.Add(2)

m := NewConcMap[string, int]()

go func() {
    defer wg.Done()
    for range 100 {
        m.Compute("hello", func(v int) int {
            return v + 1
        })
    }
}()

go func() {
    defer wg.Done()
    for range 100 {
        m.Compute("hello", func(v int) int {
            return v + 1
        })
    }
}()

wg.Wait()
fmt.Println("hello =", m.Get("hello"))
// hello = 200 (каждая горутина увеличила hello на 100)
```

Реализуйте `SetIfAbsent` и `Compute`.

