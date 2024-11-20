package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

// начало решения

var parsedDate string

// Task описывает задачу, выполненную в определенный день
type Task struct {
	Date  time.Time
	Dur   time.Duration
	Title string
}

// ParsePage разбирает страницу журнала
// и возвращает задачи, выполненные за день
func ParsePage(src string) ([]Task, error) {
	parsedDate = ""

	if src == "" {
		return []Task{}, nil
	}

	// разбиваем на строки
	lines := strings.Split(src, "\n")

	if len(lines) == 0 {
		return []Task{}, nil
	}

	date, dateError := parseDate(lines[0])

	if len(lines) == 1 {
		return []Task{}, nil
	}

	if dateError != nil {
		return []Task{}, dateError
	}

	tasks, err := parseTasks(date, lines[1:])
	if err != nil {
		return nil, err
	}
	sortTasks(tasks)
	return tasks, err
}

// parseDate разбирает дату в формате дд.мм.гггг
func parseDate(src string) (time.Time, error) {
	date, err := time.Parse("02.01.2006", src)
	parsedDate = date.Format("2006-01-02")

	return date, err
}

// parseTasks разбирает задачи из записей журнала
func parseTasks(date time.Time, lines []string) ([]Task, error) {
	if len(lines) == 0 {
		return nil, nil
	}

	tasks := make([]Task, len(lines))

	reg := regexp.MustCompile(`(\d+:\d+) - (\d+:\d+) (.+)`)

	// парсим задачи
	for i, l := range lines {
		task, err := parseTask(l, reg)

		if err != nil {
			return nil, err
		}

		task.Date = date
		tasks[i] = task
	}

	// агрегируем задачи
	tasks, _ = aggregateTasks(tasks)

	// сортируем
	sortTasks(tasks)

	return tasks, nil
}

func aggregateTasks(tasks []Task) ([]Task, error) {
	aggregatedMap := map[string]Task{}

	for _, task := range tasks {
		if t, exist := aggregatedMap[task.Title]; exist {
			task.Dur += t.Dur
		}
		aggregatedMap[task.Title] = task
	}

	result := make([]Task, len(aggregatedMap))

	i := 0
	for _, t := range aggregatedMap {
		result[i] = t
		i++
	}

	return result, nil
}

func parseTask(task string, regex *regexp.Regexp) (Task, error) {
	d := regex.FindStringSubmatch(task)

	if len(d) < 3 {
		return Task{}, fmt.Errorf("error parsing date")
	}

	timeStart, errStart := time.Parse("15:04", d[1])
	timeEnd, errEnd := time.Parse("15:04", d[2])

	if errStart != nil {
		return Task{}, errStart
	}

	if errEnd != nil {
		return Task{}, errEnd
	}

	if timeEnd.Unix() < timeStart.Unix() {
		return Task{}, fmt.Errorf("time end less than time start")
	}

	if timeEnd.Unix() == timeStart.Unix() {
		return Task{}, fmt.Errorf("time end and time start are equal")
	}

	taskName := d[3]
	taskDuration := timeEnd.Sub(timeStart)

	result := Task{
		Dur:   taskDuration,
		Title: taskName,
	}

	return result, nil
}

// sortTasks упорядочивает задачи по убыванию длительности
func sortTasks(tasks []Task) {
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Dur > tasks[j].Dur })
}

func toString(date string, tasks []Task) string {
	var result string

	if date == "" {
		return ""
	}
	result = fmt.Sprintf("Мои достижения за %s\n", date)

	for _, entry := range tasks {
		result += fmt.Sprintf("- %v: %v\n", entry.Title, entry.Dur)
	}

	result = strings.TrimSpace(result)

	return result
}

// конец решения
// ::footer

func Test(t *testing.T) {

	var tests = []struct {
		name  string
		input string
		want  string
		err   error
	}{
		{"пустое расписание", ``, ``, nil},
		{"пустое расписание 2", `15.04.2022`, `Мои достижения за 2022-04-15`, nil},
		{"пустое расписание 3", ``, ``, nil},
		{"корректное расписание", `15.04.2022
8:00 - 8:30 Завтрак
8:30 - 9:30 Оглаживание кота
9:30 - 10:00 Интернеты
10:00 - 14:00 Напряженная работа
14:00 - 14:45 Обед
14:45 - 15:00 Оглаживание кота
15:00 - 19:00 Напряженная работа
19:00 - 19:30 Интернеты
19:30 - 22:30 Безудержное веселье
22:30 - 23:00 Оглаживание кота`, `Мои достижения за 2022-04-15
- Напряженная работа: 8h0m0s
- Безудержное веселье: 3h0m0s
- Оглаживание кота: 1h45m0s
- Интернеты: 1h0m0s
- Обед: 45m0s
- Завтрак: 30m0s`, nil},
		{"расписание с ошибками", `15.04.2022
8:00 - 8:2222 Завтрак
8:30 - 9:30 Оглаживание кота
9:30 - 10:00 Интернеты`, `Мои достижения за 2022-04-15`, fmt.Errorf("parsing time \"8:2222\": extra text: \"22\"")},
	}

	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {

			got, err := ParsePage(test.input)

			if toString(parsedDate, got) != test.want {
				t.Errorf("Test %s. Got: \n%s, want: \n%s", test.name, toString(parsedDate, got), test.want)
			}

			if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", test.err) {
				t.Errorf("Test %s. Got error: %s, want: %s", test.name, fmt.Sprintf("%s", err), fmt.Sprintf("%s", test.err))
			}
		})
	}

	//entries, err := ParsePage(page)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Мои достижения за", entries[0].Date.Format("2006-01-02"))
	//for _, entry := range entries {
	//	fmt.Printf("- %v: %v\n", entry.Title, entry.Dur)
	//}

	// ожидаемый результат
	/*
		Мои достижения за 2022-04-15
		- Напряженная работа: 8h0m0s
		- Безудержное веселье: 3h0m0s
		- Оглаживание кота: 1h45m0s
		- Интернеты: 1h0m0s
		- Обед: 45m0s
		- Завтрак: 30m0s
	*/
}
