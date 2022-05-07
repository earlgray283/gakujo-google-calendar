package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/earlgray283/gakujo-google-calendar/app/crawle"
	mycalendar "github.com/earlgray283/gakujo-google-calendar/app/google-calendar-api"
	"github.com/getlantern/systray"
	calendar "google.golang.org/api/calendar/v3"
)

type App struct {
	crawler *crawle.Crawler
	srv     *calendar.Service
	Log     *log.Logger

	appLogFile     *os.File
	crawlerLogFile *os.File

	recentTaskItem     *systray.MenuItem
	recentTaskDeadLine *systray.MenuItem
	lastSyncItem       *systray.MenuItem
	unSubmittedItem    *systray.MenuItem
	syncButtonItem     *systray.MenuItem

	unSubmittedRows []*systray.MenuItem

	calendarId string
}

func NewApp(crawler *crawle.Crawler, srv *calendar.Service, logDirPath string) (*App, error) {
	appLogPath := filepath.Join(logDirPath, FmtAppLogFile)
	appFile, err := os.Create(appLogPath)
	if err != nil {
		return nil, err
	}
	logger := log.New(appFile, "", log.LstdFlags)

	crawlerLogPath := filepath.Join(logDirPath, FmtCrawlerLogFile)
	crawlerFile, err := os.Create(crawlerLogPath)
	if err != nil {
		return nil, err
	}
	crawler.Log.SetOutput(crawlerFile)

	cl, err := mycalendar.FindOrCreateCalendar("学情カレンダー", srv)
	if err != nil {
		return nil, err
	}
	clid := cl.Id

	return &App{
		crawler:        crawler,
		srv:            srv,
		Log:            logger,
		appLogFile:     appFile,
		crawlerLogFile: crawlerFile,
		calendarId:     clid,
	}, nil
}

func (a *App) Run() error {
	var err error = nil
	errC := a.crawler.Start()

	go func() {
		for {
			err = <-errC
			systray.Quit()
			return
		}
	}()

	systray.Run(a.OnReady, a.OnExit)

	return err
}
