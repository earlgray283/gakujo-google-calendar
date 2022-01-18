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
