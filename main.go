package main

import (
	"log"

	"github.com/earlgray283/gakujo-google-calendar/app"
	"github.com/earlgray283/gakujo-google-calendar/app/crawle"
	"github.com/earlgray283/gakujo-google-calendar/app/util"
)

func main() {
	config, err := app.LoadConfig(util.DefaultConfigDir())
	if err != nil {
		config, err = app.GenerateConfig()
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

	srv, err := app.NewServiceFromToken(config.Token)
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
