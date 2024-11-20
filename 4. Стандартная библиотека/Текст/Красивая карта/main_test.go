package main

import (
	"sort"
	"strconv"
	"strings"
	"testing"
)

// начало решения

// prettify возвращает отформатированное
// строковое представление карты
func prettify(m map[string]int) string {
	var keys []string
	var builder strings.Builder

	for k, _ := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// Однострочники
	if len(m) == 0 {
		builder.WriteString("{}")
		return builder.String()
	}
	if len(m) == 1 {
		builder.WriteString("{ ")
		writeKeyValue(&builder, keys[0], m[keys[0]])
		builder.WriteString(" }")
		return builder.String()
	}

	// Многострочники
	builder.WriteString("{\n")
	for _, key := range keys {
		builder.WriteString("    ")
		writeKeyValue(&builder, key, m[key])
		builder.WriteString(",\n")
	}
	builder.WriteString("}")

	return builder.String()
}

func writeKeyValue(builder *strings.Builder, key string, value int) {
	builder.WriteString(key)
	builder.WriteString(": ")
	builder.WriteString(strconv.Itoa(value))
}

// конец решения

func Test(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2, "three": 3}
	//m := map[string]int{"one": 1}
	//m := map[string]int{}
	const want = "{\n    one: 1,\n    three: 3,\n    two: 2,\n}"
	got := prettify(m)

	if got != want {
		t.Errorf("%v\ngot:\n%v\n\nwant:\n%v", m, got, want)
	}
}
