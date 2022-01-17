package app

import (
	"log"
	"testing"

	"google.golang.org/api/calendar/v3"
)

func TestInitialize(t *testing.T) {
	service, err := login()
	if err != nil {
		log.Fatal("unable to login : ", err)
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

	AddSchedule(event, calendarId, service)
}
