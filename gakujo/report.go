package gakujo

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"strconv"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
	"github.com/earlgray283/gakujo-google-calendar/gakujo/scrape"
)

func (c *Client) ReportRows(option *model.ReportSearchOption) ([]model.ReportRow, error) {
	if option.SchoolYear == 0 || option.SemesterCode.Int() == 0 {
		return nil, errors.New("some of options must be set")
	}
	page, err := c.fetchReportRowsPage(option)
	if err != nil {
		return nil, err
	}

	return scrape.ReportRows(io.NopCloser(bytes.NewReader(page)))
}

func (c *Client) ReportDetail(option *model.ReportDetailOption) (model.ReportDetail, error) {
	page, err := c.fetchReportDetail(option)
	if err != nil {
		return model.ReportDetail{}, err
	}
	return scrape.ReportDetail(io.NopCloser(bytes.NewReader(page)))
}

func (c *Client) fetchReportRowsPage(option *model.ReportSearchOption) ([]byte, error) {
	if _, err := c.fetchGeneralPurposeReportPage(); err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("schoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("semesterCode", strconv.Itoa(option.SemesterCode.Int()))
	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/report/student/searchList/search", data)
}

func (c *Client) fetchGeneralPurposeReportPage() ([]byte, error) {
	data := url.Values{}
	data.Set("headTitle", "レポート一覧")
	data.Set("menuCode", "A02")
	data.Set("nextPath", "/report/student/searchList/initialize")
	data.Set("_searchConditionDisp.accordionSearchCondition", "false")
	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/", data)
}

func (c *Client) fetchReportDetail(option *model.ReportDetailOption) ([]byte, error) {
	data := url.Values{}
	data.Set("reportId", option.ReportID)
	data.Set("listSchoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("listSubjectCode", option.ListSubjectCode)
	data.Set("listClassCode", option.ListClassCode)
	data.Set("schoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("semesterCode", strconv.Itoa(option.SemesterCode.Int()))

	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/report/student/searchList/forwardSubmitRef?submitStatusCode=01", data)
}
