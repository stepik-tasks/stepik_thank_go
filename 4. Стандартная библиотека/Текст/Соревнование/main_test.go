package main

import (
	"fmt"
	"testing"
)

// начало решения
// исходная строка может содержать любые символы;
// безопасными символами исходной строки считаются латинские буквы a-z и A-Z, цифры 0-9 и дефис -
// последовательности безопасных символов образуют «слова»;
// в результирующую строку должны попасть слова, объединенные через дефис;
// при этом буквы A-Z должны быть приведены к нижнему регистру;
// других символов в результирующей строке быть не должно.
func slugify(src string) string {
	// ..
	return src
}

// конец решения

var tests = []struct {
	title string
	want  string
}{
	{"A 100x Investment (2019)", "a-100x-investment-2019"},
}

func Test_main(t *testing.T) {
	const phrase = "A 100x Investment (2019)"
	slug := slugify(phrase)
	fmt.Println(slug)
	// a-100x-investment-2019

	for _, test := range tests {
		name := fmt.Sprintf("case(%v)", test.title)
		t.Run(name, func(t *testing.T) {
			got := slugify(test.title)
			if got != test.want {
				t.Errorf("source: %v. got %s, want %s", test.title, got, test.want)
			}
		})
	}
}
