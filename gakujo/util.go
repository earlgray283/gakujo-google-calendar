package gakujo

import (
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
)

func LatestSemesters() (int, model.SemesterCode, model.SubSemesterCode) {
	return SemestersFromTime(time.Now())
}

func SemestersFromTime(t time.Time) (int, model.SemesterCode, model.SubSemesterCode) {
	var (
		year = t.Year()
		sc   = model.EarlyPeriod
		ssc  = model.EarlyEarlyPeriod
	)
	if t.Month() <= 3 {
		year--
	}
	if 10 <= int(t.Month()) || int(t.Month()) <= 3 {
		sc = model.LaterPeriod
	}
	if 7 <= int(t.Month()) && int(t.Month()) <= 9 {
		ssc = model.EarlyLaterPeriod
	} else if 10 <= int(t.Month()) && int(t.Month()) <= 12 {
		ssc = model.LaterEarlyPeriod
	} else if int(t.Month()) <= 3 {
		ssc = model.LaterLaterPeriod
	}
	return year, sc, ssc
}
