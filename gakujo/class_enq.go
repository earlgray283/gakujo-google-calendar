package gakujo

import (
	"bytes"
	"net/url"
	"strconv"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
	"github.com/earlgray283/gakujo-google-calendar/gakujo/scrape"
)

func (c *Client) ClassEnqRows(option *model.ClassEnqSearchOption) ([]model.ClassEnqRow, error) {
	page, err := c.fetchClassEnqRowsPage(option)
	if err != nil {
		return nil, err
	}
	return scrape.ClassEnqRows(bytes.NewReader(page))
}

func (c *Client) ClassEnqDetail(option *model.ClassEnqDetailOption) (model.ClassEnqDetail, error) {
	page, err := c.fetchClassEnqDetailPage(option)
	if err != nil {
		return model.ClassEnqDetail{}, err
	}
	return scrape.ClassEnqDetail(bytes.NewReader(page))
}

func (c *Client) fetchGeneralPurposeClassEnqPage() ([]byte, error) {
	data := url.Values{}
	data.Set("headTitle", "授業評価アンケート一覧")
	data.Set("menuCode", "A05")
	data.Set("nextPath", "/classenq/student/searchList/initialize")
	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/", data)
}

func (c *Client) fetchClassEnqRowsPage(option *model.ClassEnqSearchOption) ([]byte, error) {
	if _, err := c.fetchGeneralPurposeClassEnqPage(); err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("schoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("semesterCode", strconv.Itoa(option.SemesterCode.Int()))
	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/classenq/student/searchList/search", data)
}

func (c *Client) fetchClassEnqDetailPage(option *model.ClassEnqDetailOption) ([]byte, error) {
	data := url.Values{}
	data.Set("classEnqId", option.ClassEnqID)
	data.Set("listSchoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("listSubjectCode", option.ListSubjectCode)
	data.Set("listClassCode", option.ListClassCode)
	data.Set("schoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("semesterCode", strconv.Itoa(option.SemesterCode.Int()))

	return c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/classenq/student/searchList/countingResultReference", data)
}
