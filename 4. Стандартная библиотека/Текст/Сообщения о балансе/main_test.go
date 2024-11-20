package main

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

// начало решения

var templateText = `{{.Name}}, добрый день! Ваш баланс - {{.Balance}}₽. {{if ge .Balance 100 -}} Все в порядке.{{ end -}}{{ if and (lt .Balance 100) (gt .Balance 0) -}} Пора пополнить.{{ end -}}{{ if eq .Balance 0 -}} Доступ заблокирован.{{- end}}`

// конец решения

type User struct {
	Name    string
	Balance int
}

// renderToString рендерит данные по шаблону в строку
func renderToString(tpl *template.Template, data any) string {
	var buf bytes.Buffer

	// баланс ≥ 100₽:
	// Алиса, добрый день! Ваш баланс - 1234₽. Все в порядке.

	// баланс больше 0₽, но меньше 100₽:
	// Алиса, добрый день! Ваш баланс - 77₽. Пора пополнить.

	// баланс равен 0:
	// Алиса, добрый день! Ваш баланс - 0₽. Доступ заблокирован.

	tpl.Execute(&buf, data)
	return buf.String()
}

func Test(t *testing.T) {
	var tests = []struct {
		user User
		want string
	}{
		{User{"Алиса", 101}, "Алиса, добрый день! Ваш баланс - 101₽. Все в порядке."},
		{User{"Алиса", 100}, "Алиса, добрый день! Ваш баланс - 100₽. Все в порядке."},
		{User{"Алиса", 99}, "Алиса, добрый день! Ваш баланс - 99₽. Пора пополнить."},
		{User{"Алиса", 1}, "Алиса, добрый день! Ваш баланс - 1₽. Пора пополнить."},
		{User{"Алиса", 0}, "Алиса, добрый день! Ваш баланс - 0₽. Доступ заблокирован."},
	}

	tpl := template.New("message")
	tpl = template.Must(tpl.Parse(templateText))

	for _, test := range tests {
		name := fmt.Sprintf("case(%v)", test.user)
		t.Run(name, func(t *testing.T) {
			got := renderToString(tpl, test.user)
			if got != test.want {
				t.Errorf("source: %v. got %s, want %s", test.user, got, test.want)
			}
		})
	}
}
