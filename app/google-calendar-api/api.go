package calendar

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const formatAddedAt = "2006-01-02T15:04:05+09:00"
const Scope string = calendar.CalendarScope

// コンソールにURLを表示して、コンソールにAuthCodeを貼り付けてやるやつ。テスト用。
//lint:ignore U1000 because of test
//nolint
func getLoginCodeFromStdin(URL string) (string, error) {
	fmt.Printf("Access and type logincode : \n%v\n", URL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code: %v", err)
		return "", err
	}
	return authCode, nil
}

func NewService(jsonKey []byte, token *oauth2.Token) (*calendar.Service, error) {
	ctx := context.Background()
	config, err := google.ConfigFromJSON(jsonKey, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx, token)
	return calendar.NewService(ctx, option.WithHTTPClient(client))
}

func AddSchedule(ev *calendar.Event, id string, srv *calendar.Service) error {
	isAddable, err := checkDoubleRegisted(ev.Summary, ev.End, srv, id)
	if err != nil {
		return err
	}
	if isAddable {
		_, err := srv.Events.Insert(id, ev).Do()
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadTokenFromBytes(b []byte) (*oauth2.Token, error) {
	token := &oauth2.Token{}
	if err := json.Unmarshal(b, token); err != nil {
		return nil, err
	}
	return token, nil
}

// service に title という名前のカレンダーを新規作成して、その calendar.Calendar 型を返す
func createCalendar(title string, srv *calendar.Service) (*calendar.Calendar, error) {
	newCalendar := &calendar.Calendar{
		Summary:  title,
		TimeZone: "Asia/Tokyo",
	}
	return srv.Calendars.Insert(newCalendar).Do()
}

// service に title という名前のカレンダーがあるかどうかを判定し、
// 存在すればカレンダーIDを string 型で返す
func findCalendar(title string, srv *calendar.Service) (*calendar.Calendar, error) {
	cl, err := srv.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}

	for _, item := range cl.Items {
		if item.Summary == title {
			cl, err := srv.Calendars.Get(item.Id).Do()
			if err != nil {
				return nil, err
			}
			return cl, nil
		}
	}

	return nil, nil
}

// title という名前のカレンダーがあれば calendar.Calendar 型で返す
// なければ新しく作ってそのカレンダーを calendar.Calendar 型で返す
func FindOrCreateCalendar(title string, srv *calendar.Service) (*calendar.Calendar, error) {
	cl, err := findCalendar(title, srv)
	if err != nil {
		return nil, err
	}
	if cl == nil {
		return createCalendar(title, srv)
	}
	return cl, err
}

// 予定がかぶってなかったら true を返します
func checkDoubleRegisted(eventTitle string, eventEnd *calendar.EventDateTime, srv *calendar.Service, calendarId string) (bool, error) {
	// 終了時刻の24時間前を作る
	ttt, err := time.Parse("2006-01-02T15:04:05+09:00", eventEnd.DateTime)
	if err != nil {
		return false, err
	}
	checkDateTime := ttt.Add(-24 * time.Hour).Format("2006-01-02T15:04:05+09:00")

	// 予定の取得
	events, err := srv.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(checkDateTime).MaxResults(50).OrderBy("startTime").Do()
	if err != nil {
		return false, err
	}

	for _, item := range events.Items {
		itemDateTime := item.End
		itemSummary := item.Summary
		if eventTitle == itemSummary && eventEnd.DateTime == itemDateTime.DateTime {
			return false, nil
		}
	}

	return true, nil
}

// タイトルと日時を入れると Event 型を返す
func NewGakujoEvent(title string, t time.Time) *calendar.Event {
	added := time.Now().Format(formatAddedAt)
	Event := &calendar.Event{
		Summary:     title,
		Location:    "学務情報システム",
		Description: "学情カレンダーから追加された予定です。\n学務情報システム: https://gakujo.shizuoka.ac.jp/portal/\nAdded: " + added,
		Start: &calendar.EventDateTime{
			DateTime: t.Add(-time.Hour).Format(formatAddedAt),
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			DateTime: t.Format(formatAddedAt),
			TimeZone: "Asia/Tokyo",
		},
	}
	return Event
}
