package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// начало решения

// SQLMap представляет карту, которая хранится в SQL-базе данных
type SQLMap struct {
	db         *sql.DB
	getStmt    *sql.Stmt
	setStmt    *sql.Stmt
	deleteStmt *sql.Stmt
	timeout    time.Duration
}

// NewSQLMap создает новую SQL-карту в указанной базе
func NewSQLMap(db *sql.DB) (*SQLMap, error) {
	var err error
	var getStmt, setStmt, deleteStmt *sql.Stmt

	ctx, ctxCancel := context.WithTimeout(context.Background(), 60*time.Second)

	_, err = db.ExecContext(ctx, "create table if not exists map(key text primary key, val blob)")
	if err != nil {
		ctxCancel()
		return nil, err
	}

	getStmt, err = db.PrepareContext(ctx, `select val from map where key = ?`)
	if err != nil {
		ctxCancel()
		return nil, err
	}

	setStmt, err = db.PrepareContext(ctx, `insert into map(key, val) values (?, ?) on conflict (key) do update set val = excluded.val`)
	if err != nil {
		ctxCancel()
		return nil, err
	}

	deleteStmt, err = db.PrepareContext(ctx, `delete from map where key = ?`)
	if err != nil {
		ctxCancel()
		return nil, err
	}

	ctxCancel()

	return &SQLMap{
		db,
		getStmt,
		setStmt,
		deleteStmt,
		60 * time.Second,
	}, nil
}

// SetTimeout устанавливает максимальное время выполнения
// отдельного метода карты.
func (m *SQLMap) SetTimeout(d time.Duration) {
	m.timeout = d
}

// Get возвращает значение для указанного ключа.
// Если такого ключа нет - возвращает ошибку sql.ErrNoRows.
func (m *SQLMap) Get(key string) (any, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), m.timeout)
	defer ctxCancel()

	row := m.getStmt.QueryRowContext(ctx, key)

	var e any
	err := row.Scan(&e)

	return e, err
}

// Set устанавливает значение для указанного ключа.
// Если такой ключ уже есть - затирает старое значение (это не считается ошибкой).
func (m *SQLMap) Set(key string, val any) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), m.timeout)
	defer ctxCancel()

	_, err := m.setStmt.ExecContext(ctx, []any{key, val}...)
	return err
}

// SetItems устанавливает значения указанных ключей.
func (m *SQLMap) SetItems(items map[string]any) error {
	var err error
	var tx *sql.Tx

	ctx, ctxCancel := context.WithTimeout(context.Background(), m.timeout)
	defer ctxCancel()

	tx, err = m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for k, v := range items {
		_, err = tx.Stmt(m.setStmt).ExecContext(ctx, k, v)
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
	ctx, ctxCancel := context.WithTimeout(context.Background(), m.timeout)
	defer ctxCancel()

	_, err := m.deleteStmt.ExecContext(ctx, key)
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

	m.SetTimeout(10 * time.Millisecond)

	m.Set("name", "Alice")
	a, _ := m.Get("name")

	fmt.Println(a)
}
