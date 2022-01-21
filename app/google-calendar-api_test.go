package app

import (
	//"fmt"
	//"log"
	"os"
	"testing"
	//"time"

	"github.com/earlgray283/gakujo-google-calendar/app/util"
	"google.golang.org/api/calendar/v3"
)

func TestGooglecalenderapi(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	config, err := LoadConfig(util.DefaultConfigDir())
	if err != nil {
		config, err = GenerateConfig()
		if err != nil {
			t.Fatal(err)
		}
	}

	service, err := NewServiceFromToken(config.Token)
	if err != nil {
		t.Fatal(err)
	}

	cl, err := FindOrCreateCalendar("学情カレンダー", service)
	if err != nil {
		t.Fatal(err)
	}
	calendarId := cl.Id

	// 予定を追加
	event := &calendar.Event{
		Summary:     "課題があるよ",
		Location:    "gakujo",
		Description: "はやくやれ",
		Start: &calendar.EventDateTime{
			DateTime: "2022-01-20T09:00:00+09:00",
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			DateTime: "2022-01-20T17:00:00+09:00",
			TimeZone: "Asia/Tokyo",
		},
	}

	isAbletoAdd, err := checkDoubleRegisted(event.Summary, event.End, service, calendarId)
	if err != nil {
		t.Fatal(err)
	}
	if isAbletoAdd {
		if err := AddSchedule(event, calendarId, service); err != nil {
			t.Fatal(err)
		}
	}
	//calendarId := "primary"
}
