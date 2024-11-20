package main

import (
	"fmt"
	"io"
	"strings"
)

// начало решения

// AbyssWriter пишет данные в никуда,
// но при этом считает количество записанных байт
type AbyssWriter struct {
	total int
}

func (w *AbyssWriter) Write(p []byte) (n int, err error) {
	w.total += len(p)

	return len(p), nil
}

// Total возвращает общее количество записанных байт
func (w *AbyssWriter) Total() int {
	return w.total
}

// NewAbyssWriter создает новый AbyssWriter
func NewAbyssWriter() *AbyssWriter {
	return &AbyssWriter{}
}

// конец решения

func main() {
	r := strings.NewReader("go is awesome")
	w := NewAbyssWriter()
	written, err := io.Copy(w, r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("written %d bytes\n", written)
	fmt.Println(written == int64(w.Total()))
}
