package scrape

import (
	"errors"
	"strings"
	"time"
)

func replaceAndTrim(s string) string {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	return replacer.Replace(strings.TrimSpace(s))
}

//	a wrapper of time.Parse() which supports 24:00 format
func parse2400(layout, value string) (time.Time, error) {
	parsedTime, err := time.Parse(layout, value)
	if err != nil {
		if !isHourOutErr(err) {
			return time.Time{}, err
		}
		i := strings.Index(layout, "15")
		if i == -1 {
			return time.Time{}, errors.New("stdHour 15 was not found in layout")
		}
		newValue := value[:i] + "00" + value[i+2:]
		parsedTime, err = time.Parse(layout, newValue)
		if err != nil {
			return time.Time{}, err
		}
		return parsedTime.Add(24 * time.Hour), nil
	}
	return parsedTime, nil
}

func isHourOutErr(err error) bool {
	switch err.(type) {
	case *time.ParseError:
		return strings.Contains(err.Error(), "hour")
	default:
		return false
	}
}

func toWeekday(s rune) time.Weekday {
	switch s {
	case '月':
		return time.Monday
	case '火':
		return time.Tuesday
	case '水':
		return time.Wednesday
	case '木':
		return time.Thursday
	case '金':
		return time.Friday
	case '土':
		return time.Saturday
	case '日':
		return time.Sunday
	default:
		return time.Sunday
	}
}
