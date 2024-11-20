package Високосный_год

import (
	"testing"
	"time"
)

// начало решения

func isLeapYear(year int) bool {
	t := time.Date(year, 2, 29, 0, 0, 0, 0, time.Local)
	return t.Month() == time.February
}

// конец решения

func Test(t *testing.T) {
	if !isLeapYear(2020) {
		t.Errorf("2020 is a leap year")
	}
	if isLeapYear(2022) {
		t.Errorf("2022 is NOT a leap year")
	}
}
