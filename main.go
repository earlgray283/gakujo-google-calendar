package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/earlgray283/gakujo-google-calendar/app"
	"github.com/earlgray283/gakujo-google-calendar/app/crawle"
	calendar "github.com/earlgray283/gakujo-google-calendar/app/google-calendar-api"
	"github.com/earlgray283/gakujo-google-calendar/app/util"
)

var (
	//go:embed credentials.json
	CredentialsJsonByte []byte

	logDirPath    string
	configDirPath string
	logFile       *os.File
)

func init() {
	logDirPath = util.DefaultCacheDir()
	configDirPath = util.DefaultConfigDir()
	logFile, _ = os.Create(filepath.Join(logDirPath, "log.txt"))
}

func main() {
	defer logFile.Close()

	config, err := app.LoadConfig(configDirPath)
	if err != nil {
		config, err = app.GenerateConfig(CredentialsJsonByte)
		if err != nil {
			fmt.Fprintln(logFile, err)
			return
		}
		if err := app.SaveConfig(config, configDirPath); err != nil {
			fmt.Fprintln(logFile, err)
			return
		}
	}

	crawler, err := crawle.NewCrawler(config.Username, config.Password, crawle.DefaultCrawleOption())
	if err != nil {
		fmt.Fprintln(logFile, err)
		return
	}

	srv, err := calendar.NewService(CredentialsJsonByte, config.Token)
	if err != nil {
		fmt.Fprintln(logFile, err)
		return
	}

	app, err := app.NewApp(crawler, srv, logDirPath)
	if err != nil {
		fmt.Fprintln(logFile, err)
		return
	}
	if err := app.Run(); err != nil {
		fmt.Fprintln(logFile, err)
		return
	}
}
