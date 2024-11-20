package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
)

// начало решения

type customReader struct {
	max int
	buf []byte
}

func (c *customReader) Read(p []byte) (n int, err error) {
	n, err = rand.Read(c.buf)

	copy(p, c.buf)

	return len(c.buf), io.EOF
}

// RandomReader создает читателя, который возвращает случайные байты,
// но не более max штук
func RandomReader(max int) io.Reader {
	return bufio.NewReaderSize(&customReader{
		max: max,
		buf: make([]byte, max),
	}, max)
}

// конец решения

func main() {
	rnd := RandomReader(5)
	rd := bufio.NewReader(rnd)
	for {
		b, err := rd.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d ", b)
	}
	fmt.Println()
	// 1 148 253 194 250
	// (значения могут отличаться)
}
