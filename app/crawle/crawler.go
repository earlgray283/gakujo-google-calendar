package crawle

import (
	"log"
	"sync"
	"time"

	"github.com/earlgray283/gakujo-google-calendar/app/util"
	"github.com/earlgray283/gakujo-google-calendar/gakujo"
	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
	"github.com/go-co-op/gocron"
)

type Crawler struct {
	gc  *GakujoClient
	opt *CrawleOption
	Log *log.Logger

	username string
	password string

	Minitest *MinitestStorage
	Report   *ReportStorage
	Classenq *ClassEnqStorage
}

type GakujoClient struct {
	c *gakujo.Client
	sync.Mutex
}

func NewCrawler(username, password string, opt *CrawleOption) (*Crawler, error) {
	gc := gakujo.NewClient()
	err := util.DoWithRetry(func() error {
		return gc.Login(username, password)
	}, 5, 5*time.Second)
	if err != nil {
		return nil, err
	}
	wgc := &GakujoClient{c: gc}
	rs := &ReportStorage{rows: []model.ReportRow{}}
	ms := &MinitestStorage{rows: []model.MinitestRow{}}
	ces := &ClassEnqStorage{rows: []model.ClassEnqRow{}}
	return &Crawler{wgc, opt, log.Default(), username, password, ms, rs, ces}, nil
}

func (c *Crawler) Start() chan error {
	s := gocron.NewScheduler(time.Local)
	errc := make(chan error)
	_, _ = s.Every(c.opt.MinitestInterval).Do(func() {
		if err := c.CrawleMinitestRows(c.opt.RetryCount); err != nil {
			c.Log.Println(err)
			errc <- err
		}
	})
	_, _ = s.Every(c.opt.ReportInterval).Do(func() {
		if err := c.CrawleReportRows(c.opt.RetryCount); err != nil {
			c.Log.Println(err)
			errc <- err
		}
	})
	_, _ = s.Every(c.opt.ClassenqInterval).Do(func() {
		if err := c.CrawleClassEnqRows(c.opt.RetryCount); err != nil {
			c.Log.Println(err)
			errc <- err
		}
	})
	// 30分ごとにセッションを回復する
	_, _ = s.Every(30).Minutes().Do(func() {
		c.gc.Lock()
		defer c.gc.Unlock()
		c.Log.Println("updating session")
		if _, err := c.gc.c.LatestClassEnqRows(); err != nil {
			if err := c.gc.c.Login(c.username, c.password); err != nil {
				errc <- err
			}
		}
		c.Log.Println("succeed in updating session")
	})

	s.StartAsync()

	go func() {
		for {
			<-errc
			s.Stop()
			s.Clear()
			return
		}
	}()

	return errc
}
