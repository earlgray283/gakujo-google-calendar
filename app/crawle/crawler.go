package crawle

import (
	"log"
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo"
	"github.com/go-co-op/gocron"
)

type Crawler struct {
	gc  *gakujo.Client
	opt *CrawleOption
	Log *log.Logger

	username string
	password string

	minitest *MinitestStorage
	report   *ReportStorage
	classenq *ClassEnqStorage
}

func NewCrawler(username, password string, opt *CrawleOption) (*Crawler, error) {
	gc := gakujo.NewClient()
	if err := gc.Login(username, password); err != nil {
		return nil, err
	}
	return &Crawler{gc, opt, log.Default(), username, password, nil, nil, nil}, nil
}

func (c *Crawler) Start() chan error {
	s := gocron.NewScheduler(time.Local)
	errc := make(chan error)
	tasks := []struct {
		f        func(int) error
		interval time.Duration
	}{
		{c.crawleMinitestRows, c.opt.MinitestInterval},
		{c.crawleReportRows, c.opt.ReportInterval},
		{c.crawleClassEnqRows, c.opt.ClassenqInterval},
	}
	for _, task := range tasks {
		s.Every(uint64(task.interval.Hours())).Hours().Do(func() {
			if err := task.f(c.opt.RetryCount); err != nil {
				c.Log.Println(err)
				errc <- err
			}
		})
	}
	// 30分ごとにセッションを回復する
	s.Every(30).Minutes().Do(func() {
		if _, err := c.gc.LatestClassEnqRows(); err != nil {
			if err := c.gc.Login(c.username, c.password); err != nil {
				errc <- err
			}
		}
	})

	s.StartAsync()

	go func() {
		for {
			select {
			case <-errc:
				s.Stop()
				s.Clear()
				return
			}
		}
	}()

	return errc
}
