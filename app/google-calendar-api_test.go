package app

import (
	"os"
	"testing"
	"log"
	"fmt"

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

	/* テストの度に実行されてカレンダーが増殖するので一旦コメントアウト！
	// 新しいカレンダーを作成
	newCalendar := &calendar.Calendar{
		Summary: "学情カレンダー",
	}
	createdCalendar, err := service.Calendars.Insert(newCalendar).Do()
	if err != nil {
		t.Fatal(err)
	}

	// カレンダーIDを取得
	calendarId := createdCalendar.Id
	fmt.Println(calendarId)
	*/


	// カレンダー一覧を取得
    cl, err := service.CalendarList.List().Do()
    if err != nil {
        log.Fatalf("Error CalendarList.List(): %v", err)
        return
    }

	// カレンダー一覧の名前とIDをプリント
    fmt.Printf("--- Your calendars ---\n")
    for _, item := range cl.Items {
        fmt.Println(item.Summary + " " + item.Id)
    }

	// 予定を追加
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
	}
	calendarId := "primary"
	if err := AddSchedule(event, calendarId, service); err != nil {
		t.Fatal(err)
	}
}
