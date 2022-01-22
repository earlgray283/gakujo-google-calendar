package calendar

import (
	//"fmt"
	//"log"
	"os"
	"path/filepath"
	"testing"

	//"time"

	"github.com/earlgray283/gakujo-google-calendar/app/util"
	"google.golang.org/api/calendar/v3"
)

func TestGooglecalenderapi(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	configDir := util.DefaultConfigDir()
	b, err := os.ReadFile(filepath.Join(configDir, "token.json"))
	if err != nil {
		t.Fatal(err)
	}
	token, err := LoadTokenFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	credentialsJsonByte, _ := os.ReadFile("../credentials.json")

	service, err := NewService(credentialsJsonByte, token)
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

	isAddable, err := checkDoubleRegisted(event.Summary, event.End, service, calendarId)
	if err != nil {
		t.Fatal(err)
	}
	if isAddable {
		if err := AddSchedule(event, calendarId, service); err != nil {
			t.Fatal(err)
		}
	}
	//calendarId := "primary"
}
