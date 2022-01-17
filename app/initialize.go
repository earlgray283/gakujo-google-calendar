package app

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"
	"regexp"
)

type User struct {
	Username  string
	Password  string
	Logincode string
}

func GetUserInfoFromBrowser(url string) (User, error) {
	UserInfo := User{}
	mux := http.NewServeMux()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("template/auth.html")
		if err != nil {
			log.Fatalf("template error: %v", err)
		}
		if err := t.Execute(rw, struct {
			URL string
		}{
			URL: url,
		}); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
	})

	mux.HandleFunc("/redirect", func(rw http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("template/redirect.html")
		if err != nil {
			log.Fatalf("template error: %v", err)
		}
		if err := t.Execute(rw, struct {
			URL string
		}{
			URL: url,
		}); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
	})

	mux.HandleFunc("/regist", func(rw http.ResponseWriter, r *http.Request) {
		if m, _ := regexp.MatchString("^[a-zA-Z0-9]+$", r.FormValue("username")); !m {
			http.Redirect(rw, r, "/redirect", 301)
			return
		}
		UserInfo.Username = r.FormValue("username")
		
		if m, _ := regexp.MatchString("^[a-zA-Z0-9]+$", r.FormValue("password")); !m {
			http.Redirect(rw, r, "/redirect", 301)
			return
		}
		UserInfo.Password = r.FormValue("password")
		
		if m, _ := regexp.MatchString("^[a-zA-Z0-9!-/:-@¥[-`{-~]+$", r.FormValue("logincode")); !m {
			http.Redirect(rw, r, "/redirect", 301)
			return
		}
		UserInfo.Logincode = r.FormValue("logincode")
		
		if len(UserInfo.Username) == 0 || len(UserInfo.Password) == 0 || len(UserInfo.Logincode) == 0 {
			http.Error(rw, "username, password, logincode must not be empty", http.StatusBadRequest)
			return
		}

		// ここでシャットダウンするお
		log.Println("Server shutdown")
		_ = srv.Shutdown(ctx)

		//fmt.Println(UserInfo.Username, UserInfo.Password, UserInfo.Token)
	})

	log.Println("Listening on port http://localhost:8080")

	srv.Handler = mux

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return UserInfo, err
	}

	return UserInfo, nil

}
