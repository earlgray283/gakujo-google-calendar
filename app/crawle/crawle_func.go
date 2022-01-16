package crawle

import (
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
)

func (c *Crawler) crawleMinitestRows(retryCount int) error {
	var (
		err  error
		rows []model.MinitestRow
	)
	c.minitest.Lock()
	defer c.minitest.Unlock()
	for i := 0; i < retryCount; i++ {
		rows, err = c.gc.LatestMinitestRows()
		if err == nil {
			break
		}
		c.Log.Println("error occurred. retry after 5s...:", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return err
	}
	c.minitest.Rows = rows
	c.minitest.UpdatedAt = time.Now()
	return nil
}

func (c *Crawler) crawleReportRows(retryCount int) error {
	var (
		err  error
		rows []model.ReportRow
	)
	c.report.Lock()
	defer c.report.Unlock()
	for i := 0; i < retryCount; i++ {
		rows, err = c.gc.LatestReportRows()
		if err == nil {
			break
		}
		c.Log.Println("error occurred. retry after 5s...:", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return err
	}
	c.report.Rows = rows
	c.report.UpdatedAt = time.Now()
	return nil
}

func (c *Crawler) crawleClassEnqRows(retryCount int) error {
	var (
		err  error
		rows []model.ClassEnqRow
	)
	c.classenq.Lock()
	defer c.classenq.Unlock()
	for i := 0; i < retryCount; i++ {
		rows, err = c.gc.LatestClassEnqRows()
		if err == nil {
			break
		}
		c.Log.Println("error occurred. retry after 5s...:", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return err
	}
	c.classenq.Rows = rows
	c.classenq.UpdatedAt = time.Now()
	return nil
}
