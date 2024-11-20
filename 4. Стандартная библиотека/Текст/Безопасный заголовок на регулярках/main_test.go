package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// начало решения

// последовательность допустимых символов
var wordRE = regexp.MustCompile(`[a-z0-9\-]+`)

// slugify возвращает "безопасный" вариант заголовока:
// только латиница, цифры и дефис
// Требования:
//
//	исходная строка может содержать любые символы;
//	безопасными символами исходной строки считаются латинские буквы a-z и A-Z, цифры 0-9 и дефис -
//	последовательности безопасных символов образуют «слова»;
//	в результирующую строку должны попасть слова, объединенные через дефис;
//	при этом буквы A-Z должны быть приведены к нижнему регистру;
//	других символов в результирующей строке быть не должно.
func slugify(src string) string {
	words := wordRE.FindAllString(strings.ToLower(src), -1)
	return strings.Join(words, "-")
}

// конец решения

func Test(t *testing.T) {
	var tests = []struct {
		str  string
		want string
	}{
		{"Go Is Awesome!", "go-is-awesome"},
		{"!!123Go Is Awesome!", "123go-is-awesome"},
		{"Go - Is - Awesome", "go---is---awesome"},
		{"!Attention, attention!", "attention-attention"},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case(%s)", test.str)
		t.Run(name, func(t *testing.T) {
			got := slugify(test.str)
			if got != test.want {
				t.Errorf("source: %s. got %s, want %s", test.str, got, test.want)
			}
		})
	}
}
