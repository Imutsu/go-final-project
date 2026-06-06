package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, date.Location())

	return d.After(n)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid d format")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("invalid day interval")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid y format")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	case "w":
		return nextWeekDate(now, date, parts)

	case "m":
		return "", errors.New("unsupported month format")

	default:
		return "", errors.New("unsupported repeat format")
	}
}

func nextWeekDate(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid week format")
	}

	var days [8]bool

	items := strings.Split(parts[1], ",")

	for _, item := range items {
		n, err := strconv.Atoi(item)
		if err != nil || n < 1 || n > 7 {
			return "", errors.New("invalid weekday")
		}

		days[n] = true
	}

	for {
		date = date.AddDate(0, 0, 1)

		wd := int(date.Weekday())

		if wd == 0 {
			wd = 7
		}

		if days[wd] && afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
