package scrape

import (
	"fmt"
	"strings"
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
)

// return beginDate, endDate, error
func parsePeriod(periodText string) (time.Time, time.Time, error) {
	var beginText1, beginText2, endText1, endText2 string
	fmt.Sscanf(periodText, "%s %s ～ %s %s", &beginText1, &beginText2, &endText1, &endText2)
	beginDate, err := parse2400("2006/01/02 15:04", fmt.Sprintf("%s %s", beginText1, beginText2))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endDate, err := parse2400("2006/01/02 15:04", fmt.Sprintf("%s %s", endText1, endText2))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return beginDate, endDate, nil
}

func parseCourseNameFormat(s string) (string, []model.CourseDate, error) {
	s = strings.TrimSpace(s)
	elems := strings.Split(s, "\n")
	if len(elems) != 2 {
		return "", nil, fmt.Errorf("invalid course name format: ===%s===", s)
	}
	for i := range elems {
		elems[i] = strings.TrimSpace(elems[i])
	}

	courseDates := []model.CourseDate{}
	for _, plainCourseDate := range strings.Split(elems[1], ",") {
		courseDate, err := parseCourseDateFormat(plainCourseDate)
		if err != nil {
			return "", nil, err
		}
		courseDates = append(courseDates, *courseDate)
	}

	return elems[0], courseDates, nil
}

func parseCourseDateFormat(s string) (*model.CourseDate, error) {
	var (
		semester    string
		weekday     rune
		jigen1      int
		jigen2      int
		subSemester string
		other       string
	)

	elms := strings.Split(s, "/")
	if len(elms) != 2 {
		return nil, fmt.Errorf("invalid course date format: ===%s===", s)
	}
	fmt.Sscanf(elms[0], "%s", &semester)

	// 前期/水5・6
	if _, err := fmt.Sscanf(elms[1], "%c%d・%d", &weekday, &jigen1, &jigen2); err != nil {
		if _, err := fmt.Sscanf(elms[1], "%c%d・%d(%s)", &weekday, &jigen1, &jigen2, &subSemester); err != nil {
			if _, err := fmt.Sscanf(elms[1], "%s", &other); err != nil {
				return nil, fmt.Errorf("invalid course date format: ===%s===", s)
			}
		}
	}
	return &model.CourseDate{
		SemesterCode:    model.ToSemesterCode(semester),
		Weekday:         toWeekday(weekday),
		Jigen1:          jigen1,
		Jigen2:          jigen2,
		SubSemesterCode: model.ToSubSemesterCode(subSemester),
		Other:           other,
	}, nil
}
