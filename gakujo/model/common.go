package model

import "time"

type SubSemesterCode string

const (
	None             = SubSemesterCode("")
	EarlyEarlyPeriod = SubSemesterCode("前期前半")
	EarlyLaterPeriod = SubSemesterCode("前期後半")
	LaterEarlyPeriod = SubSemesterCode("後期前半")
	LaterLaterPeriod = SubSemesterCode("後期後半")
)

type SemesterCode string

const (
	SemesterCodeNone = SemesterCode("")
	EarlyPeriod      = SemesterCode("前期")
	LaterPeriod      = SemesterCode("後期")
)

type CourseDate struct {
	SemesterCode    SemesterCode
	Weekday         time.Weekday
	Jigen1          int
	Jigen2          int
	SubSemesterCode SubSemesterCode
	Other           string
}

func ToSemesterCode(s string) SemesterCode {
	switch s {
	case "前期":
		return EarlyPeriod
	case "後期":
		return LaterPeriod
	default:
		return SemesterCodeNone
	}
}

func (sc SemesterCode) Int() int {
	switch sc {
	case EarlyPeriod:
		return 1
	case LaterPeriod:
		return 2
	case SemesterCodeNone:
		return 0
	default:
		return 0
	}
}

func ToSubSemesterCode(s string) SubSemesterCode {
	switch s {
	case "前期前半":
		return EarlyEarlyPeriod
	case "前期後半":
		return EarlyLaterPeriod
	case "後期前半":
		return LaterEarlyPeriod
	case "後期後半":
		return LaterLaterPeriod
	default:
		return ""
	}
}
