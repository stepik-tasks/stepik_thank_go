**Карта в базе**

Реализуйте карто-подобный тип SQLMap, который хранит данные в таблице базы. Ключи карты — только строки, а значения — произвольного типа.

Интерфейс:

```go
// Get возвращает значение для указанного ключа.
// Если такого ключа нет - возвращает ошибку sql.ErrNoRows.
Get(key string) (any, error)

// Set устанавливает значение для указанного ключа.
// Если такой ключ уже есть - затирает старое значение (это не считается ошибкой).
Set(key string, val any) error

// Delete удаляет запись карты с указанным ключом.
// Если такого ключа нет - ничего не делает (это не считается ошибкой).
Delete(key string) error

```

Пример использования:

```go
db, err := sql.Open("sqlite3", "map.db")
// ...

m, err := NewSQLMap(db)
if err != nil {
    panic(err)
}

m.Set("name", "Alice")
m.Set("age", 42)

name, err := m.Get("name")
fmt.Printf("name = %v, err = %v\n", name, err)
// name = Alice, err = <nil>

age, err := m.Get("age")
fmt.Printf("age = %v, err = %v\n", age, err)
// age = 42, err = <nil>

m.Set("name", "Bob")
name, err = m.Get("name")
fmt.Printf("name = %v, err = %v\n", name, err)
// name = Bob, err = <nil>

m.Delete("name")
name, err = m.Get("name")
fmt.Printf("name = %v, err = %v\n", name, err)
// name = <nil>, err = sql: no rows in result set
```

**SQL-запросы для SQLite**

Храните данные в таблице map. Инициализируйте ее таким запросом:

```sql
create table if not exists map(key text primary key, val blob)
```

Не меняйте название таблицы, а также названия и типы столбцов — они используются в тестах.

Вот запросы, которые должен выполнять каждый из методов:

- Get
```sql
select val from map where key = ?
```

- Set
```sql
insert into map(key, val) values (?, ?)
on conflict (key) do update set val = excluded.val
```

Это так называемый upsert (insert or update). Он пытается добавить новую запись в таблицу, но если строка с таким ключом уже есть — обновляет ее, затирая старое значение val. Синтаксис upsert-запросов отличается в разных СУБД. У SQLite такой.

- Delete
```sql
delete from map where key = ?
```

**SQL-запросы для MySQL**
```sql
-- create
create table if not exists map(`key` varchar(255) primary key, val blob)

-- get
select val from map where `key` = ?

-- set
replace into map(`key`, val) values (?, ?)

-- delete
delete from map where `key` = ?
```