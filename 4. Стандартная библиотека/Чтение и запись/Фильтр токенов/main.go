package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// TokenReader начитывает токены из источника
type TokenReader interface {
	// ReadToken считывает очередной токен
	// Если токенов больше нет, возвращает ошибку io.EOF
	ReadToken() (string, error)
}

// TokenWriter записывает токены в приемник
type TokenWriter interface {
	// WriteToken записывает очередной токен
	WriteToken(s string) error
}

// начало решения

type TReader struct {
	reader *bufio.Scanner
}

func (receiver TReader) ReadToken() (string, error) {
	if !receiver.reader.Scan() {
		return "", io.EOF
	}

	return receiver.reader.Text(), receiver.reader.Err()
}

type TWriter struct {
	words int
}

func (receiver TWriter) WriteToken(s string) error {
	receiver.words++
	return nil
}

func (receiver TWriter) Words() int {
	return receiver.words
}

// FilterTokens читает все токены из src и записывает в dst тех,
// кто проходит проверку predicate
func FilterTokens(dst TokenWriter, src TokenReader, predicate func(s string) bool) (int, error) {
	total := 0

	for {
		token, err := src.ReadToken()

		if err == io.EOF {
			break
		}

		if err != nil {
			return total, err
		}

		if predicate(token) {
			err = dst.WriteToken(token)
			if err != nil {
				return total, err
			}
			total++
		}
	}

	return total, nil
}

func NewWordReader(src string) TReader {
	scanner := bufio.NewScanner(strings.NewReader(src))
	scanner.Split(bufio.ScanWords)

	return TReader{
		reader: scanner,
	}
}

func NewWordWriter() TWriter {
	return TWriter{
		words: 0,
	}
}

// конец решения

func main() {
	// Для проверки придется создать конкретные типы,
	// которые реализуют интерфейсы TokenReader и TokenWriter.

	// Ниже для примера используются NewWordReader и NewWordWriter,
	// но вы можете сделать любые на свое усмотрение.

	r := NewWordReader("go is awesome")
	w := NewWordWriter()
	predicate := func(s string) bool {
		return s != "is"
	}
	n, err := FilterTokens(w, r, predicate)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d tokens: %v\n", n, w.Words())
	// 2 tokens: [go awesome]
}
