package app

import (
	"context"
	_ "embed"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

//go:embed credentials.json
var CredentialsJsonByte []byte

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

func NewServiceFromToken(token *oauth2.Token) (*calendar.Service, error) {
	ctx := context.Background()
	config, err := google.ConfigFromJSON(CredentialsJsonByte, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx, token)
	return calendar.NewService(ctx, option.WithHTTPClient(client))
}

func AddSchedule(ev *calendar.Event, id string, srv *calendar.Service) error {
	_, err := srv.Events.Insert(id, ev).Do()
	if err != nil {
		return err
	}

	return nil
}

func createNewCalendar(title string, srv *calendar.Service) (*calendar.Calendar, error) {
	// service に title という名前のカレンダーを新規作成して、その calendar.Calendar 型を返す
	newCalendar := &calendar.Calendar{
		Summary: title,
	}
	createdCalendar, err := srv.Calendars.Insert(newCalendar).Do()
	if err != nil {
		return nil, err
	}

	return createdCalendar, nil
}

func exploreCalendar(title string, srv *calendar.Service) (*calendar.Calendar, error) {
	// service に title という名前のカレンダーがあるかどうかを判定し、
	// 存在すればカレンダーIDを string 型で返す
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

func FindCalendar(title string, srv *calendar.Service) (*calendar.Calendar, error) {
	// title という名前のカレンダーがあれば calendar.Calendar 型で返す
	// なければ新しく作ってそのカレンダーを calendar.Calendar 型で返す
	cl, err := exploreCalendar(title, srv)
	if err != nil {
		return nil, err
	} else {
		if cl == nil {
			// カレンダーが存在しないので作成
			newCalendar, err := createNewCalendar(title, srv)
			if err != nil {
				return nil, err
			}
			return newCalendar, nil
		} else {
			return cl, err
		}
	}
}