package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// начало решения

// SQLMap представляет карту, которая хранится в SQL-базе данных
type SQLMap struct {
	db         *sql.DB
	getStmt    *sql.Stmt
	setStmt    *sql.Stmt
	deleteStmt *sql.Stmt
}

// NewSQLMap создает новую SQL-карту в указанной базе
func NewSQLMap(db *sql.DB) (*SQLMap, error) {
	var err error
	var getStmt, setStmt, deleteStmt *sql.Stmt

	_, err = db.Exec("create table if not exists map(key text primary key, val blob)")
	if err != nil {
		return nil, err
	}

	getStmt, err = db.Prepare(`select val from map where key = ?`)
	if err != nil {
		return nil, err
	}

	setStmt, err = db.Prepare(`insert into map(key, val) values (?, ?) on conflict (key) do update set val = excluded.val`)
	if err != nil {
		return nil, err
	}

	deleteStmt, err = db.Prepare(`delete from map where key = ?`)
	if err != nil {
		return nil, err
	}

	return &SQLMap{db, getStmt, setStmt, deleteStmt}, nil
}

// Get возвращает значение для указанного ключа.
// Если такого ключа нет - возвращает ошибку sql.ErrNoRows.
func (m *SQLMap) Get(key string) (any, error) {
	row := m.getStmt.QueryRow(key)

	var e any
	err := row.Scan(&e)

	return e, err
}

// Set устанавливает значение для указанного ключа.
// Если такой ключ уже есть - затирает старое значение (это не считается ошибкой).
func (m *SQLMap) Set(key string, val any) error {
	_, err := m.setStmt.Exec([]any{key, val}...)
	return err
}

// SetItems устанавливает значения указанных ключей.
func (m *SQLMap) SetItems(items map[string]any) error {
	var err error
	var tx *sql.Tx

	tx, err = m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for k, v := range items {
		_, err = tx.Stmt(m.setStmt).Exec(k, v)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

// Delete удаляет запись карты с указанным ключом.
// Если такого ключа нет - ничего не делает (это не считается ошибкой).
func (m *SQLMap) Delete(key string) error {
	_, err := m.deleteStmt.Exec(key)
	return err
}

// Close освобождает ресурсы, занятые картой в базе.
func (m *SQLMap) Close() error {
	var err error

	err = m.getStmt.Close()
	if err != nil {
		return err
	}

	err = m.setStmt.Close()
	if err != nil {
		return err
	}

	err = m.deleteStmt.Close()
	if err != nil {
		return err
	}

	return nil
}

// конец решения

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	m, err := NewSQLMap(db)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	m.Set("name", "Alice")

	items := map[string]any{
		"name": "Bob",
		"age":  42,
	}
	m.SetItems(items)

	name, err := m.Get("name")
	fmt.Printf("name = %v, err = %v\n", name, err)
	// name = Bob, err = <nil>

	age, err := m.Get("age")
	fmt.Printf("age = %v, err = %v\n", age, err)
	// age = 42, err = <nil>
}
