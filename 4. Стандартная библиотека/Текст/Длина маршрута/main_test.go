package main

import (
	"strconv"
	"strings"
	"testing"
)

// начало решения

// calcDistance возвращает общую длину маршрута в метрах
func calcDistance(directions []string) int {
	totalDistance := 0
	for _, section := range directions {

		w := strings.Fields(section)
		for _, d := range w {
			// метры
			if d[len(d)-1] == 'm' && !strings.Contains(d, "km") {
				t := strings.ReplaceAll(d, "m", "")
				d, _ := strconv.Atoi(t)
				totalDistance += d
			}
			// километры
			if d[len(d)-1] == 'm' && strings.Contains(d, "km") {
				t := strings.ReplaceAll(d, "km", "")
				d, _ := strconv.ParseFloat(t, 64)
				totalDistance += int(d * 1000)
			}
		}

	}

	return totalDistance
}

// конец решения

func Test(t *testing.T) {
	directions := []string{
		"straight 1.6km",
		//"100m to intersection",
		//"turn right",
		//"straight 300m",
		//"enter motorway",
		//"straight 5km",
		//"exit motorway",
		//"500m straight",
		//"turn sharp left",
		//"continue 100m to destination",
	}
	const want = 6000
	got := calcDistance(directions)
	if got != want {
		t.Errorf("%v: got %v, want %v", directions, got, want)
	}
}
