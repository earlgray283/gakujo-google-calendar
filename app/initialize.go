package app

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/earlgray283/gakujo-google-calendar/assets"
	"github.com/skratchdot/open-golang/open"
)

type AuthFormInfo struct {
	Username  string
	Password  string
	Logincode string
}

func GetAuthInfoFromBrowser(url string) (*AuthFormInfo, error) {
	authFormInfo := AuthFormInfo{}
	errC := make(chan error)
	sigC := make(chan struct{})

	t, err := template.ParseFS(assets.HtmlAuth, "html/auth.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if err := t.Execute(rw, struct {
			URL string
		}{
			URL: url,
		}); err != nil {
			http.Error(rw, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/redirect", func(rw http.ResponseWriter, r *http.Request) {
		if err := t.Execute(rw, struct {
			URL string
		}{
			URL: url,
		}); err != nil {
			http.Error(rw, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/regist", func(rw http.ResponseWriter, r *http.Request) {
		if m, _ := regexp.MatchString("^[a-zA-Z0-9]+$", r.FormValue("username")); !m {
			http.Redirect(rw, r, "/redirect", http.StatusMovedPermanently)
			return
		}
		if m, _ := regexp.MatchString("^[a-zA-Z0-9]+$", r.FormValue("password")); !m {
			http.Redirect(rw, r, "/redirect", http.StatusMovedPermanently)
			return
		}
		if m, _ := regexp.MatchString("^[a-zA-Z0-9!-/:-@¥[-`{-~]+$", r.FormValue("logincode")); !m {
			http.Redirect(rw, r, "/redirect", http.StatusMovedPermanently)
			return
		}

		authFormInfo.Username = r.FormValue("username")
		authFormInfo.Password = r.FormValue("password")
		authFormInfo.Logincode = r.FormValue("logincode")

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("Registration succeed! You can close this page."))

		sigC <- struct{}{}
	})

	if err := open.Run("http://localhost:8080"); err != nil {
		return nil, err
	}
	srv.Handler = mux
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case err := <-errC:
			return nil, err
		case <-sigC:
			log.Println("Server shutdown")
			_ = srv.Shutdown(context.Background())
			return &authFormInfo, nil
		}
	}
}
