package main

import (
	"fmt"
	"strings"
	"testing"
)

// начало решения

// 	Требования:
//
//	исходная строка может содержать любые символы;
//	безопасными символами исходной строки считаются латинские буквы a-z и A-Z, цифры 0-9 и дефис -
//	последовательности безопасных символов образуют «слова»;
//	в результирующую строку должны попасть слова, объединенные через дефис;
//	при этом буквы A-Z должны быть приведены к нижнему регистру;
//	других символов в результирующей строке быть не должно.

/**
ASCII
abcdefghijklmnopqrstuvwxyz: 97-122 | 65-90
1234567890: 48-57
-: 45

Допустимые диапазоны: 45, 48-57, 65-90, 97-122

A - 65
a - 97

B - 66
b - 98

D - 68
d - 100

пробел - 32
*/

// slugify возвращает "безопасный" вариант заголовока:
func slugify(src string) string {
	var source []byte
	var builder strings.Builder
	var ch byte
	var i int

	source = []byte(src)
	builder = strings.Builder{}
	builder.Grow(len(source))

	// цикл по всем символам
	for i, ch = range source {
		// допустимые символы
		if ch == 45 || (ch >= 48 && ch <= 57) || (ch >= 65 && ch <= 90) || ch >= 97 && ch <= 122 {
			// верхний регистр преобразуется в нижний
			if ch >= 65 && ch <= 90 {
				builder.WriteByte(ch + 32)
			} else {
				builder.WriteByte(ch)
			}
		} else {
			// последовательности безопасных символов образуют «слова»;
			// следующий символ должен быть безопасным
			if i > 0 &&
				i < len(source)-1 &&
				builder.Len() < i+2 &&
				(source[i+1] == 45 ||
					(source[i+1] >= 48 && source[i+1] <= 57) ||
					(source[i+1] >= 65 && source[i+1] <= 90) ||
					source[i+1] >= 97 && source[i+1] <= 122) {
				builder.WriteByte(45)
			}
		}
	}

	return builder.String()
}

// конец решения

var tests = []struct {
	source string
	want   string
}{
	{
		"!Attention, attention!",
		"attention-attention",
	},
	{
		"We haven't killed 90% of all plankton",
		"we-haven-t-killed-90-of-all-plankton",
	},
	{
		"Carbon Language: An experimental successor to C++",
		"carbon-language-an-experimental-successor-to-c",
	},
	{
		"Hello, 中国!",
		"hello",
	},
	{
		"Tz6t5bx S9zne Fw-6i Giv0f F894; Tp-.",
		"tz6t5bx-s9zne-fw-6i-giv0f-f894-tp-",
	},
	{
		"Zkaab41ov Lk- Yde0c3xc Wo9e12n17 F-5h-ysbv Yzxn& R9uhm236h",
		"zkaab41ov-lk--yde0c3xc-wo9e12n17-f-5h-ysbv-yzxn-r9uhm236h",
	},
	{
		"Go Talks: \"Cuddle: an App Engine Demo\"",
		"go-talks-cuddle-an-app-engine-demo",
	},
}

func Test_main(t *testing.T) {
	for _, test := range tests {
		name := fmt.Sprintf("case(%v)", test.source)
		t.Run(name, func(t *testing.T) {
			got := slugify(test.source)
			if got != test.want {
				t.Errorf("source: %v. got %s, want %s", test.source, got, test.want)
			}
		})
	}
}

func Benchmark_main(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, test := range tests {
			slugify(test.source)
		}
	}
}
