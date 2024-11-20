package main

import (
	"fmt"
	"strings"
	"testing"
)

// начало решения

// slugify возвращает "безопасный" вариант заголовока:
// только латиница, цифры и дефис
const dict = "abcdefghijklmnopqrstuvwxyz1234567890-"

func slugify(src string) string {
	var result []string

	words := strings.Fields(src)

	// вырезаем лишние символы
	//var clr string
	for _, word := range words {
		if w := clearWord(word); w != "" {
			result = append(result, clearWord(word))
		} else {
			result = append(result, "-")
		}
	}

	return strings.Trim(strings.Join(result, "-"), "-")
}

func clearWord(word string) string {
	// сюда будем собирать символы
	result := make([]string, len(word))

	// книжнему регистру
	word = strings.ToLower(word)

	// обрезаем с краев лишнее
	word = strings.Trim(word, "!@#$%^&*()_+=:,")
	last := ""

	for _, ch := range word {
		if strings.Contains(dict, string(ch)) {
			result = append(result, string(ch))
			last = string(ch)
		} else {
			if last != "-" {
				result = append(result, "-")
				last = string(ch)
			}
		}
	}

	return strings.Trim(strings.Join(result, ""), "-")
}

// конец решения

func main() {
	// JSON-RPC: a tale of interfaces   ->   "json-rpc--a-tale-of-interfaces", want "json-rpc-a-tale-of-interfaces"
	// Hello, 中国!: ->  "hello-", want "hello"
	// Go's New Brand
	// Debugging Go code (a status report)
	//r := slugify("Go - Is - Awesome") // -> go---is---awesome
	//r := slugify("Debugging Go code (a status report)") // -> debugging-go-code-a-status-report
	//r := slugify("Hello, 中国!:") // -> debugging-go-code-a-status-report
	//r := slugify("JSON-RPC: a tale of interfaces") // -> json-rpc-a-tale-of-interfaces
	//r := slugify("Arrays, slices (and strings): The mechanics of 'append'") // -> arrays-slices-and-strings-the-mechanics-of-append
	r := slugify("Go_Is_Awesome")

	fmt.Println(r)
}

func Test(t *testing.T) {
	const phrase = "Go Is Awesome!"
	const want = "go-is-awesome"
	got := slugify(phrase)
	if got != want {
		t.Errorf("%s: got %#v, want %#v", phrase, got, want)
	}
}
