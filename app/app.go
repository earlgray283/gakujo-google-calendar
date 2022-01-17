package app

import (
	"log"

	"github.com/earlgray283/gakujo-google-calendar/app/crawle"
	"google.golang.org/api/calendar/v3"
)

type App struct {
	crawler *crawle.Crawler
	srv     *calendar.Service
	Log     *log.Logger
}

func NewApp(crawler *crawle.Crawler, srv *calendar.Service) *App {
	return &App{
		crawler: crawler,
		srv:     srv,
		Log:     log.Default(),
	}
}
