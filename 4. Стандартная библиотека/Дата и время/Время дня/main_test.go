package main

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// начало решения

// TimeOfDay описывает время в пределах одного дня
type TimeOfDay struct {
	hour     int
	minute   int
	second   int
	location *time.Location
}

// Hour возвращает часы в пределах дня
func (t TimeOfDay) Hour() int {
	return t.hour
}

// Minute возвращает минуты в пределах часа
func (t TimeOfDay) Minute() int {
	return t.minute
}

// Second возвращает секунды в пределах минуты
func (t TimeOfDay) Second() int {
	return t.second
}

// String возвращает строковое представление времени
// в формате чч:мм:сс TZ (например, 12:34:56 UTC)
func (t TimeOfDay) String() string {
	return fmt.Sprintf("%.2d:%.2d:%.2d %s", t.hour, t.minute, t.second, t.location.String())
}

// Equal сравнивает одно время с другим.
// Если у t и other разные локации - возвращает false.
func (t TimeOfDay) Equal(other TimeOfDay) bool {
	if t.location.String() != other.location.String() {
		return false
	}

	return t.String() == other.String()
}

// Before возвращает true, если время t предшествует other.
// Если у t и other разные локации - возвращает ошибку.
func (t TimeOfDay) Before(other TimeOfDay) (bool, error) {
	if t.location.String() != other.location.String() {
		return false, errors.New("different time zones")
	}

	// переводим в кол-во секунд с начала дня 00:00:00
	return (t.hour*60*60 + t.minute*60 + t.second) < (other.hour*60*60 + other.minute*60 + other.second), nil
}

// After возвращает true, если время t идет после other.
// Если у t и other разные локации - возвращает ошибку.
func (t TimeOfDay) After(other TimeOfDay) (bool, error) {
	if t.location.String() != other.location.String() {
		return false, errors.New("different time zones")
	}

	// переводим в кол-во секунд с начала дня 00:00:00
	return (t.hour*60*60 + t.minute*60 + t.second) > (other.hour*60*60 + other.minute*60 + other.second), nil
}

// MakeTimeOfDay создает время в пределах дня
func MakeTimeOfDay(hour, min, sec int, loc *time.Location) TimeOfDay {
	return TimeOfDay{hour, min, sec, loc}
}

// конец решения

func Test(t *testing.T) {
	t1 := MakeTimeOfDay(17, 45, 22, time.UTC)
	t2 := MakeTimeOfDay(20, 3, 4, time.UTC)

	if t1.Equal(t2) {
		t.Errorf("%v should not be equal to %v", t1, t2)
	}

	before, _ := t1.Before(t2)
	if !before {
		t.Errorf("%v should be before %v", t1, t2)
	}

	after, _ := t1.After(t2)
	if after {
		t.Errorf("%v should NOT be after %v", t1, t2)
	}
}
