package main

import (
	"bufio"
	"fmt"
	"os"
)

// начало решения

// readLines возвращает все строки из указанного файла
func readLines(name string) ([]string, error) {
	descriptor, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(descriptor)

	var result []string

	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// конец решения

func main() {
	lines, err := readLines("/Users/a.krizhanovsky/webhome/stepik/st/stepik_thank_go/4. Стандартная библиотека/Чтение и запись/Строки из файла (снова)/test.txt")

	if err != nil {
		panic(err)
	}
	for idx, line := range lines {
		fmt.Printf("%d: %s\n", idx+1, line)
	}
}
