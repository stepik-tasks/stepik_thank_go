package main

import (
	"fmt"
	mathrand "math/rand"
	"os"
	"path/filepath"
)

// алфавит планеты Нибиру
const alphabet = "aeiourtnsl"

// Census реализует перепись населения.
// Записи о рептилоидах хранятся в каталоге census, в отдельных файлах,
// по одному файлу на каждую букву алфавита.
// В каждом файле перечислены рептилоиды, чьи имена начинаются
// на соответствующую букву, по одному рептилоиду на строку.
type Census struct {
	dir      string
	dirEntry os.DirEntry
}

// Count возвращает общее количество переписанных рептилоидов.
func (c *Census) Count() int {
	return 0
}

// Add записывает сведения о рептилоиде.
func (c *Census) Add(name string) {
	l := string(name[0])

	fileName := fmt.Sprintf("%s/%s.txt", c.dir, l)

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Cannot open %s to write reptiloid", fileName))
	}

	_, err = f.WriteString(fmt.Sprintf("%s\n", name))
	if err != nil {
		fmt.Println(err)
		panic("Cannot write reptiloid")
	}
}

// Close закрывает файлы, использованные переписью.
func (c *Census) Close() {
}

// NewCensus создает новую перепись и пустые файлы
// для будущих записей о населении.
func NewCensus() *Census {
	dir, err := filepath.Abs("./census")

	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(dir, 0777)
	if err != nil {
		return nil
	}

	return &Census{dir: dir}
}

// ┌─────────────────────────────────┐
// │ не меняйте код ниже этой строки │
// └─────────────────────────────────┘

var rand = mathrand.New(mathrand.NewSource(0))

// randomName возвращает имя очередного рептилоида.
func randomName(n int) string {
	chars := make([]byte, n)
	for i := range chars {
		chars[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(chars)
}

func main() {
	census := NewCensus()
	defer census.Close()
	for i := 0; i < 1024; i++ {
		//for i := 0; i < 10; i++ {
		reptoid := randomName(5)
		census.Add(reptoid)
	}
	fmt.Println(census.Count())
}
