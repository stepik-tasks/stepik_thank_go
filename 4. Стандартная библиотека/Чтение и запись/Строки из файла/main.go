package main

import (
	"fmt"
	"os"
	"strings"
)

// начало решения

// readLines возвращает все строки из указанного файла
func readLines(name string) ([]string, error) {
	content, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	resultSlice := strings.Split(string(content), "\n")

	var result []string

	for _, line := range resultSlice {
		if line != "" {
			result = append(result, line)
		}
	}

	return result, nil
}

// конец решения

func main() {
	lines, err := readLines("/Users/a.krizhanovsky/webhome/stepik/st/stepik_thank_go/4. Стандартная библиотека/Чтение и запись/test.txt")
	if err != nil {
		panic(err)
	}
	for idx, line := range lines {
		fmt.Printf("%d: %s\n", idx+1, line)
	}
}
