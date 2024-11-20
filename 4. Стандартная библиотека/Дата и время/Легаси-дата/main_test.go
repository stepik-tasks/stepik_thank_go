package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// начало решения

var parseError = fmt.Errorf("error while parse timestamp")

// asLegacyDate преобразует время в легаси-дату
func asLegacyDate(t time.Time) string {
	result := strings.TrimRight(fmt.Sprintf("%d.%.9d", t.Unix(), t.Nanosecond()), "0")

	if result[len(result)-1:] == "." {
		result = fmt.Sprintf("%s0", result)
	}

	return result
}

// parseLegacyDate преобразует легаси-дату во время.
// Возвращает ошибку, если легаси-дата некорректная.
func parseLegacyDate(d string) (time.Time, error) {
	reg := regexp.MustCompile(`(\d+)\.(\d+)`)

	matches := reg.FindStringSubmatch(d)

	if len(matches) < 3 {
		return time.Time{}, parseError
	}

	// парсим секунды
	seconds, e := strconv.Atoi(matches[1])
	if e != nil {
		return time.Time{}, parseError
	}

	// наносекунды
	nanoSecondsSlice := []string{"0", "0", "0", "0", "0", "0", "0", "0", "0"}

	for i, c := range matches[2] {
		nanoSecondsSlice[i] = string(c)
	}

	nanoSeconds, e := strconv.Atoi(strings.Join(nanoSecondsSlice, ""))
	if e != nil {
		return time.Time{}, parseError
	}

	result := time.Unix(int64(float64(seconds)), int64(nanoSeconds))

	return result, nil
}

// конец решения

func Test_asLegacyDate(t *testing.T) {

	//    main.go:80: 1970-01-01 01:00:00.000000123 +0000 UTC: got 3600.123, want 3600.000000123
	//    main.go:80: 2022-05-24 14:45:22.951205 +0000 UTC: got 1653403522.951205000, want 1653403522.951205
	//    main.go:80: 1970-01-01 01:00:00.000123456 +0000 UTC: got 3600.123456, want 3600.000123456
	//    main.go:80: 2022-05-24 14:45:22.951 +0000 UTC: got 1653403522.951000000, want 1653403522.951
	//    main.go:80: 1970-01-01 01:00:00.123456 +0000 UTC: got 3600.123456000, want 3600.123456

	// 	  main.go:120: 3600.123: got 1970-01-01 01:00:00.000000123 +0000 UTC, want 1970-01-01 01:00:00.123 +0000 UTC
	samples := map[time.Time]string{
		time.Date(1970, 1, 1, 1, 0, 0, 123, time.UTC):       "3600.000000123",
		time.Date(1970, 1, 1, 1, 0, 0, 1, time.UTC):         "3600.000000001",
		time.Date(1970, 1, 1, 1, 0, 0, 123456789, time.UTC): "3600.123456789",
		time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC):         "3600.0",
		time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC):         "0.0",
	}
	for src, want := range samples {
		got := asLegacyDate(src)
		if got != want {
			t.Fatalf("%v: got %v, want %v", src, got, want)
		}
	}
}

func Test_parseLegacyDate(t *testing.T) {

	//main.go:122: 1653403522.951000: got 2022-05-24 14:45:22.000951 +0000 UTC, want 2022-05-24 14:45:22.951 +0000 UTC
	//main.go:122: 1653403522.951205: got 2022-05-24 14:45:22.000951205 +0000 UTC, want 2022-05-24 14:45:22.951205 +0000 UTC
	//main.go:122: 1653403522.951: got 2022-05-24 14:45:22.000000951 +0000 UTC, want 2022-05-24 14:45:22.951 +0000 UTC
	//main.go:122: 3600.123: got 1970-01-01 01:00:00.000000123 +0000 UTC, want 1970-01-01 01:00:00.123 +0000 UTC
	//main.go:122: 3600.123456: got 1970-01-01 01:00:00.000123456 +0000 UTC, want 1970-01-01 01:00:00.123456 +0000 UTC

	samples := map[string]time.Time{
		"1653403522.951000": time.Date(2022, 05, 24, 14, 45, 22, 951, time.UTC),

		//"1653403522.951000": time.Date(2022, 5, 24, 14, 45, 22, 951, time.UTC),
		//"1653403522.951205": time.Date(2022, 5, 24, 14, 45, 22, 951205, time.UTC),
		//"3600.123456789":    time.Date(1970, 1, 1, 1, 0, 0, 123456789, time.UTC),
		//"3600.0":            time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC),
		//"0.0":               time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		//"1.123456789":       time.Date(1970, 1, 1, 0, 0, 1, 123456789, time.UTC),
	}
	for src, want := range samples {
		got, err := parseLegacyDate(src)
		if err != nil {
			t.Fatalf("%v: unexpected error", src)
		}
		if !got.Equal(want) {
			t.Fatalf("%v: got %v, want %v", src, got, want)
		}
	}
}
