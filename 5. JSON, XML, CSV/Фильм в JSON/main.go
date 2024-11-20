package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// начало решения

// Duration описывает продолжительность фильма
type Duration time.Duration

// MarshalJSON
// Правила кодирования продолжительности фильма:
// строка вида XhYm, где X — количество часов, а Y — количество минут (2h15m);
// если часов 0, то они опускаются (не 0h45m, а 45m);
// если минут 0, то они опускаются (не 2h0m, а 2h).
func (d Duration) MarshalJSON() ([]byte, error) {
	var b strings.Builder
	var hours int
	var minutes int

	hours = int(time.Duration(d).Hours())
	minutes = int(time.Duration(d).Minutes()) - (int(hours) * 60)

	// 23h31m
	b.WriteString("\"")

	if hours > 0 {
		b.WriteString(fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		b.WriteString(fmt.Sprintf("%dm", minutes))
	}

	b.WriteString("\"")

	return []byte(b.String()), nil
}

// Rating описывает рейтинг фильма
type Rating int

// MarshalJSON
// Правила кодирования рейтинга:
//
// рейтинг принимает значения от 0 до 5;
// выходная строка всегда состоит из 5 звезд;
// из них заполнены X звезд, где X — значение рейтинга (3 = ★★★☆☆).
func (r Rating) MarshalJSON() ([]byte, error) {
	var b strings.Builder

	b.WriteString("\"")
	for i := 1; i <= 5; i++ {
		if i <= int(r) {
			b.WriteString("★")
		} else {
			b.WriteString("☆")
		}
	}
	b.WriteString("\"")

	return []byte(b.String()), nil
}

// Movie описывает фильм
type Movie struct {
	Title    string
	Year     int
	Director string
	Genres   []string
	Duration Duration
	Rating   Rating
}

// MarshalMovies кодирует фильмы в JSON.
//   - если indent = 0 - использует json.Marshal
//   - если indent > 0 - использует json.MarshalIndent
//     с отступом в указанное количество пробелов.
func MarshalMovies(indent int, movies ...Movie) (string, error) {
	var err error
	var r []byte

	if indent == 0 {
		r, err = json.Marshal(movies)
	} else {
		r, err = json.MarshalIndent(movies, "", "    ")
	}
	return string(r), err
}

// конец решения

func main() {
	m1 := Movie{
		Title:    "Interstellar",
		Year:     2014,
		Director: "Christopher Nolan",
		Genres:   []string{"Adventure", "Drama", "Science Fiction"},
		Duration: Duration(2*time.Hour + 49*time.Minute),
		Rating:   5,
	}
	m2 := Movie{
		Title:    "Sully",
		Year:     2016,
		Director: "Clint Eastwood",
		Genres:   []string{"Drama", "History"},
		Duration: Duration(time.Hour + 36*time.Minute),
		Rating:   4,
	}

	s, err := MarshalMovies(4, m1, m2)
	fmt.Println(err)
	// nil
	fmt.Println(s)
	/*
		[
		    {
		        "Title": "Interstellar",
		        "Year": 2014,
		        "Director": "Christopher Nolan",
		        "Genres": [
		            "Adventure",
		            "Drama",
		            "Science Fiction"
		        ],
		        "Duration": "2h49m",
		        "Rating": "★★★★★"
		    },
		    {
		        "Title": "Sully",
		        "Year": 2016,
		        "Director": "Clint Eastwood",
		        "Genres": [
		            "Drama",
		            "History"
		        ],
		        "Duration": "1h36m",
		        "Rating": "★★★★☆"
		    }
		]
	*/
}
