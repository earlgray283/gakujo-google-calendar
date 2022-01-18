package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

const (
	FmtCrawlerLogFile = "crawler_log_2006_01_02_15_04_05.txt"
	FmtAppLogFile     = "app_log_2006_01_02_15_04_05.txt"
)

type Config struct {
	Username string
	Password string
	Token    *oauth2.Token
}

func LoadConfig(configDir string) (*Config, error) {
	envMap, err := godotenv.Read(filepath.Join(configDir, ".account"))
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Join(configDir, "token.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, err
	}
	return &Config{
		Username: envMap["GAKUJO_USERNAME"],
		Password: envMap["GAKUJO_PASSWORD"],
		Token:    tok,
	}, nil
}

func SaveConfig(config *Config, configDir string) error {
	tokenPath := filepath.Join(configDir, "token.json")
	accPath := filepath.Join(configDir, ".account")
	f1, err := os.OpenFile(accPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f1.Close()
	_, _ = f1.WriteString("GAKUJO_USERNAME" + "=" + config.Username + "\n")
	_, _ = f1.WriteString("GAKUJO_PASSWORD" + "=" + config.Password + "\n")
	f2, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f2.Close()
	return json.NewEncoder(f2).Encode(config.Token)
}

func GenerateConfig() (*Config, error) {
	ctx := context.Background()

	fmt.Println(string(CredentialsJsonByte))

	//スコープの設定
	config, err := google.ConfigFromJSON(CredentialsJsonByte, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}

	// authentication url を生成
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// AuthCodeを入力させる
	authFormInfo, err := GetAuthInfoFromBrowser(authURL)
	if err != nil {
		return nil, err
	}

	token, err := config.Exchange(ctx, authFormInfo.Logincode)
	if err != nil {
		return nil, err
	}

	return &Config{
		Username: authFormInfo.Username,
		Password: authFormInfo.Password,
		Token:    token,
	}, nil
}
