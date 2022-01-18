package app

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/earlgray283/gakujo-google-calendar/app/crawle"
	"github.com/earlgray283/gakujo-google-calendar/app/util"
	"github.com/getlantern/systray"
	"google.golang.org/api/calendar/v3"
)

type App struct {
	crawler *crawle.Crawler
	srv     *calendar.Service
	Log     *log.Logger

	appLogFile     *os.File
	crawlerLogFile *os.File
}

func NewApp(crawler *crawle.Crawler, srv *calendar.Service) (*App, error) {
	appLogPath := filepath.Join(util.DefaultConfigDir(), time.Now().Format(FmtAppLogFile))
	appFile, err := os.Create(appLogPath)
	if err != nil {
		return nil, err
	}
	logger := log.New(appFile, "", log.LstdFlags)

	crawlerLogPath := filepath.Join(util.DefaultConfigDir(), time.Now().Format(FmtCrawlerLogFile))
	crawlerFile, err := os.Create(crawlerLogPath)
	if err != nil {
		return nil, err
	}
	crawler.Log.SetOutput(crawlerFile)

	return &App{
		crawler:        crawler,
		srv:            srv,
		Log:            logger,
		appLogFile:     appFile,
		crawlerLogFile: crawlerFile,
	}, nil
}

func (a *App) Run() error {
	errC := a.crawler.Start()

	systray.Run(a.OnReady, a.OnExit)

	for {
		err := <-errC
		systray.Quit()
		return err
	}
}
