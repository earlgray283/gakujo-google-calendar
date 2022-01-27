package crawle

import (
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
)

func (c *Crawler) CrawleMinitestRows(retryCount int) error {
	var (
		err  error
		rows []model.MinitestRow
	)
	c.Minitest.Lock()
	defer c.Minitest.Unlock()
	c.gc.Lock()
	defer c.gc.Unlock()
	c.Log.Println("start crawling minitest")
	for i := 0; i < retryCount; i++ {
		rows, err = c.gc.c.LatestMinitestRows()
		if err == nil {
			break
		}
		c.Log.Println("error occurred. retry after 5s...:", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return err
	}
	c.Minitest.rows = rows
	c.Minitest.updatedAt = time.Now()
	c.Log.Printf("succeed in crawling minitest(%d rows)", len(rows))
	return nil
}

func (c *Crawler) CrawleReportRows(retryCount int) error {
	var (
		err  error
		rows []model.ReportRow
	)
	c.Report.Lock()
	defer c.Report.Unlock()
	c.gc.Lock()
	defer c.gc.Unlock()
	c.Log.Println("start crawling report")
	for i := 0; i < retryCount; i++ {
		rows, err = c.gc.c.LatestReportRows()
		if err == nil {
			break
		}
		c.Log.Println("error occurred. retry after 5s...:", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return err
	}
	c.Report.rows = rows
	c.Report.updatedAt = time.Now()
	c.Log.Printf("succeed in crawling report(%d rows)", len(rows))
	return nil
}

func (c *Crawler) CrawleClassEnqRows(retryCount int) error {
	var (
		err  error
		rows []model.ClassEnqRow
	)
	c.Classenq.Lock()
	defer c.Classenq.Unlock()
	c.gc.Lock()
	defer c.gc.Unlock()
	c.Log.Println("start crawling classenq")
	for i := 0; i < retryCount; i++ {
		rows, err = c.gc.c.LatestClassEnqRows()
		if err == nil {
			break
		}
		c.Log.Println("error occurred. retry after 5s...:", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return err
	}
	c.Classenq.rows = rows
	c.Classenq.updatedAt = time.Now()
	c.Log.Printf("succeed in crawling classenq(%d rows)", len(rows))
	return nil
}
