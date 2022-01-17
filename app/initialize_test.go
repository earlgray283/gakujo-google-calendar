package app

import (
	"fmt"
	"log"
	"testing"
)

func TestGetUserInfoFromBrowser(t *testing.T) {
	url := "https://github.com/earlgray283/gakujo-google-calendar/pull/13/files"
	UserInfo, err := GetUserInfoFromBrowser(url)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(UserInfo.Username)
}

// go test -run ^TestInitialize$ github.com/earlgray283/gakujo-google-calendar/app -v -count=1
