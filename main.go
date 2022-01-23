package main

import (
	_ "embed"
	"log"

	"github.com/earlgray283/gakujo-google-calendar/app"
	"github.com/earlgray283/gakujo-google-calendar/app/crawle"
	calendar "github.com/earlgray283/gakujo-google-calendar/app/google-calendar-api"
	"github.com/earlgray283/gakujo-google-calendar/app/util"
)

//go:embed credentials.json
var CredentialsJsonByte []byte

func main() {
	config, err := app.LoadConfig(util.DefaultConfigDir())
	if err != nil {
		config, err = app.GenerateConfig(CredentialsJsonByte)
		if err != nil {
			log.Fatal(err)
		}
		if err := app.SaveConfig(config, util.DefaultConfigDir()); err != nil {
			log.Fatal(err)
		}
	}

	crawler, err := crawle.NewCrawler(config.Username, config.Password, crawle.DefaultCrawleOption())
	if err != nil {
		log.Fatal(err)
	}

	srv, err := calendar.NewService(CredentialsJsonByte, config.Token)
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.NewApp(crawler, srv)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
