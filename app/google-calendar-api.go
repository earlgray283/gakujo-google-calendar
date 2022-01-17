package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func makeAuthURL(config *oauth2.Config) string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

func getTokenWithAuthCode(aCode string, config *oauth2.Config) (*oauth2.Token, error) {
	authCode := aCode

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

//コンソールにURLを表示して、コンソールにAuthCodeを貼り付けてやるやつ。テスト用。
func GetLoginCodeFromStdin(URL string) (string, error) {
	fmt.Printf("Access and type logincode : \n%v\n", URL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code: %v", err)
		return "", err
	}
	return authCode, nil
}

// 2回目以降のログインに使う
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
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

func login(credentialFile string, tokenFile string) (*calendar.Service, error) {
	//Google Cloud Platform 上のアプリケーションにアクセスするための情報
	ctx := context.Background()
	b, err := os.ReadFile(credentialFile)
	if err != nil {
		return nil, err
	}

	//スコープの設定
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}

	// URLとる
	authURL := makeAuthURL(config)

	// AuthCodeを入力させる
	UserInfo, err := GetUserInfoFromBrowser(authURL)
	if err != nil {
		return nil, err
	}
	authCode := UserInfo.Logincode

	// Tokenを取る
	tok, err := getTokenWithAuthCode(authCode, config)
	if err != nil {
		return nil, err
	}

	// Tokenを保存する
	tokFile := tokenFile
	err = saveToken(tokFile, tok)
	if err != nil {
		return nil, err
	}

	/* 2回目以降のログインは↓のようにファイルからトークンを読む
	tok, _ := tokenFromFile("token.json")
	(エラーハンドリングは省略)
	*/
	client := config.Client(context.Background(), tok)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, err
}

func AddSchedule(ev *calendar.Event, id string, srv *calendar.Service) error {
	_, err := srv.Events.Insert(id, ev).Do()
	if err != nil {
		return err
	}

	return nil
}
