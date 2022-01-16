package app

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
)

type User struct {
	Username string
	Password string
	Token    string
}

func GetUserInfoFromBrowser() User{
	UserInfo := User{}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("template/auth.html")
		if err != nil {
			log.Fatalf("template error: %v", err)
		}
		if err := t.Execute(rw, struct {		
			URL string
		}{
			URL:"https://github.com/earlgray283/gakujo-google-calendar",
		}); err != nil {
		log.Printf("failed to execute template: %v", err)
		}
	})

	http.HandleFunc("/resist", func(rw http.ResponseWriter, r *http.Request){
		UserInfo.Username = r.FormValue("username")
		UserInfo.Password = r.FormValue("password")
		UserInfo.Token = r.FormValue("token")
		if len(UserInfo.Username) == 0 || len(UserInfo.Password) == 0 || len(UserInfo.Token) == 0 {
    		http.Error(rw, "username, password, token must not be empty", http.StatusBadRequest)
    	return
		}
		fmt.Println(UserInfo.Username, UserInfo.Password, UserInfo.Token)
	})

	log.Println("Listening on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	return UserInfo
}