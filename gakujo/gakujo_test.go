package gakujo

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/model"
	"github.com/joho/godotenv"
)

var (
	c *Client
)

func init() {
	if err := godotenv.Load("./../.env"); err != nil {
		if len(os.Getenv("GAKUJO_USERNAME")) == 0 || len(os.Getenv("GAKUJO_PASSWORD")) == 0 {
			log.Fatal(err)
		}
	}
	username, password := os.Getenv("GAKUJO_USERNAME"), os.Getenv("GAKUJO_PASSWORD")
	log.Println(username, password)
	c = NewClient()
	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	inc := NewClient()
	username, password := os.Getenv("GAKUJO_USERNAME"), os.Getenv("GAKUJO_PASSWORD")
	log.Println(username, password)
	if err := inc.Login(username, password); err != nil {
		t.Fatal(err)
	}
}

func TestReportRows(t *testing.T) {
	option := model.ReportSearchOption{
		SchoolYear:   2021,
		SemesterCode: model.LaterPeriod,
	}
	reportRows, err := c.ReportRows(&option)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range reportRows {
		fmt.Println(row)
	}
}

func TestMinitestRows(t *testing.T) {
	option := model.MinitestSearchOption{
		SchoolYear:   2021,
		SemesterCode: model.LaterPeriod,
	}
	minitestRows, err := c.MinitestRows(&option)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range minitestRows {
		fmt.Println(row)
	}
}
