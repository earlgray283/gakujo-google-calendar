package gakujo

import (
	"bytes"
	"io"
	"net/url"
	"strconv"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
	"github.com/earlgray283/gakujo-google-calendar/gakujo/scrape"
)

func (c *Client) MinitestRows(option *model.MinitestSearchOption) ([]model.MinitestRow, error) {
	page, err := c.fetchMinitestRowsPage(option)
	if err != nil {
		return nil, err
	}
	return scrape.MinitestRows(io.NopCloser(bytes.NewReader(page)))
}

func (c *Client) MinitestDetail(option *model.MinitestDetailOption) (model.MinitestDetail, error) {
	page, err := c.fetchMinitestDetailPage(option)
	if err != nil {
		return model.MinitestDetail{}, err
	}
	return scrape.MinitestDetail(io.NopCloser(bytes.NewReader(page)))
}

func (c *Client) fetchMinitestRowsPage(option *model.MinitestSearchOption) ([]byte, error) {
	if _, err := c.fetchGeneralPurposeMinitestRowsPage(); err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("schoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("semesterCode", strconv.Itoa(option.SemesterCode.Int()))
	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/test/student/searchList/search", data)
}

func (c *Client) fetchGeneralPurposeMinitestRowsPage() ([]byte, error) {
	data := url.Values{}
	data.Set("headTitle", "小テスト一覧")
	data.Set("menuCode", "A03")
	data.Set("nextPath", "/test/student/searchList/initialize")
	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/", data)
}

func (c *Client) fetchMinitestDetailPage(option *model.MinitestDetailOption) ([]byte, error) {
	data := url.Values{}
	data.Set("testId", option.TestID)
	data.Set("listSchoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("listSubjectCode", option.ListSubjectCode)
	data.Set("listClassCode", option.ListClassCode)
	data.Set("schoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("semesterCode", strconv.Itoa(option.SemesterCode.Int()))

	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/test/student/searchList/forwardSubmitRef", data)
}
