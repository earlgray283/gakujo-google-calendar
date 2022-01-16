package scrape

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
)

func MinitestRows(r io.Reader) ([]model.MinitestRow, error) {
	var err error
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	rows := []model.MinitestRow{}
	selection := doc.Find("#searchList > tbody > tr")
	selection.EachWithBreak(func(i int, s *goquery.Selection) bool {
		courseName, courseDates, inerr := parseCourseNameFormat(s.Find("td:nth-child(1)").Text())
		if inerr != nil {
			err = inerr
			return false
		}

		title := replaceAndTrim(s.Find("td:nth-child(2)").Text())

		jsText, exists := s.Find("td:nth-child(2) > a").Attr("onclick")
		if !exists {
			err = errors.New(`attr "onClick" not found`)
			return false
		}
		minitestMetadata, inerr := parseMinitestJSargument(jsText)
		if inerr != nil {
			err = inerr
			return false
		}

		statusText := replaceAndTrim(s.Find("td:nth-child(3)").Text())
		status := model.ToStatus(statusText)

		periodText := replaceAndTrim(s.Find("td:nth-child(4)").Text())
		beginDate, endDate, inerr := parsePeriod(periodText)
		if inerr != nil {
			err = inerr
			return false
		}

		submitStatusText := replaceAndTrim(s.Find("td:nth-child(5)").Text())
		submitStatus := model.ToSubmitStatus(submitStatusText)

		format := replaceAndTrim(s.Find("td:nth-child(6)").Text())
		rows = append(rows, model.MinitestRow{
			CourseName:   courseName,
			CourseDates:  courseDates,
			Title:        title,
			Status:       status,
			BeginDate:    beginDate,
			EndDate:      endDate,
			SubmitStatus: submitStatus,
			Format:       format,
			TaskMetadata: minitestMetadata,
		})
		return true
	})
	return rows, err
}

func MinitestDetail(r io.Reader) (model.MinitestDetail, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return model.MinitestDetail{}, err
	}
	selection := doc.Find("#area > table > tbody")
	title := replaceAndTrim(selection.Find("tr:nth-child(1) > td").Text())
	periodText := replaceAndTrim(selection.Find("tr:nth-child(2) > td").Text())
	beginDate, endDate, err := parsePeriod(periodText)
	if err != nil {
		return model.MinitestDetail{}, err
	}
	numText := replaceAndTrim(selection.Find("tr:nth-child(3) > td").Text())
	var num int
	fmt.Sscanf(numText, "%d 問", &num)
	evaluationMethod := replaceAndTrim(selection.Find("tr:nth-child(4) > td").Text())
	description := replaceAndTrim(selection.Find("tr:nth-child(5) > td").Text())
	description = strings.Join(strings.Split(description, "<br/>"), "\n")
	transMatter := replaceAndTrim(selection.Find("tr:nth-child(7) > td").Text())
	var minitestHtml string
	doc.Find("#area > div:nth-child(4) > table").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i < 2 {
			return true
		}
		html, inerr := s.Html()
		if err != nil {
			err = inerr
			return false
		}
		minitestHtml += html
		return true
	})
	if err != nil {
		return model.MinitestDetail{}, err
	}
	return model.MinitestDetail{
		Title:            title,
		BeginDate:        beginDate,
		EndDate:          endDate,
		Num:              num,
		EvaluationMethod: evaluationMethod,
		Description:      description,
		TransMatter:      transMatter,
		MinitestHtml:     minitestHtml,
	}, nil
}

func parseMinitestJSargument(jsArgument string) (model.TaskMetadata, error) {
	tokens := strings.Split(jsArgument[11:len(jsArgument)-2], ",")
	for i, token := range tokens {
		newToken := replaceAndTrim(token)
		tokens[i] = newToken[1 : len(newToken)-1]
	}
	if len(tokens) != 6 {
		return model.TaskMetadata{}, errors.New("too few tokens")
	}

	year, err := strconv.Atoi(tokens[3])
	if err != nil {
		return model.TaskMetadata{}, err
	}
	return model.TaskMetadata{
		ID:               tokens[1],
		SubmitStatusCode: tokens[2],
		SchoolYear:       year,
		SubjectCode:      tokens[4],
		ClassCode:        tokens[5],
	}, nil
}
