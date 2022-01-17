package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hoge/app"
	"os"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func makeAuthURL(config *oauth2.Config) string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

func getTokenWithAuthCode(aCode string, config *oauth2.Config) (*oauth2.Token, error) {
	authCode := aCode
	if len(authCode) == 0 {
		fmt.Println("Unable to read authorization code")
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token from web: %v", err)
		return nil, err
	}
	return tok, nil
}

//コンソールにURLを表示して、コンソールにAuthCodeを貼り付けてやるやつ。テスト用。
func PrintAuthURL (URL string) (string, error) {
	fmt.Printf("Access and type logincode : \n%v\n", URL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code: %v", err)
		return "", err
	}
	return authCode, nil
} 

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Unable to cache oauth token: %v", err)
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

func login() (*calendar.Service, error) {
	//Google Cloud Platform 上のアプリケーションにアクセスするための情報
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v", err)
		return nil, err
	}

	//スコープの設定
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	// URLとる
	authURL := makeAuthURL(config)

	// AuthCodeを入力させる
	UserInfo := app.GetUserInfoFromBrowser(authURL)
	authCode := UserInfo.Logincode

	// Tokenを取る
	tok, err := getTokenWithAuthCode(authCode, config)
	if err != nil {
		fmt.Printf("Unable to retrive token : %v", err)
		return nil, err
	}

	// Tokenを保存する
	tokFile := "token.json"
	saveToken(tokFile, tok)
	
	/* 2回目以降のログインは↓のようにファイルからトークンを読む
	tok, _ := tokenFromFile("token.json")
	(エラーハンドリングは省略)
	*/
	client := config.Client(context.Background(), tok)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Print("Unable to retrieve Calendar client: %v", err)
		return nil, err
	}

	return srv, err
}

func AddSchedule(ev *calendar.Event, id string, srv *calendar.Service) error {
	_, err := srv.Events.Insert(id, ev).Do()
	if err != nil {
		fmt.Printf("Unable to create event. %v\n", err)
		return err
	}

	return nil
}