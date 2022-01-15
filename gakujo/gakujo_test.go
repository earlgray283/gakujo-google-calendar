package gakujo

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("./../.env"); err != nil {
		log.Println(err)
	}
}

func TestLogin(t *testing.T) {
	inc := NewClient()
	username, password := os.Getenv("GAKUJO_USERNAME"), os.Getenv("GAKUJO_PASSWORD")
	log.Println(username, password)
	if err := inc.Login(username, password); err != nil {
		t.Fatal(err)
	}
}
