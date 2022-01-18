package app

import (
	"os"
	"testing"

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

	//ここから予定の追加
	// Refer to the Go quickstart on how to setup the environment:
	// https://developers.google.com/calendar/quickstart/go
	// Change the scope to calendar.CalendarScope and delete any stored credentials.
	event := &calendar.Event{
		Summary:     "課題があるよ",
		Location:    "gakujo",
		Description: "はやくやれ",
		Start: &calendar.EventDateTime{
			DateTime: "2022-01-20T09:00:00+09:00",
			TimeZone: "Asia/Tokyo",
			//TimeZone: "America/Los_Angeles",
		},
		End: &calendar.EventDateTime{
			DateTime: "2022-01-20T17:00:00+09:00",
			TimeZone: "Asia/Tokyo",
			//TimeZone: "America/Los_Angeles",
		},
		//Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=2"},
		/*
		   Attendees: []*calendar.EventAttendee{
		   &calendar.EventAttendee{Email:"lpage@example.com"},
		   &calendar.EventAttendee{Email:"sbrin@example.com"},
		   },
		*/
	}

	calendarId := "primary"

	if err := AddSchedule(event, calendarId, service); err != nil {
		t.Fatal(err)
	}
}
