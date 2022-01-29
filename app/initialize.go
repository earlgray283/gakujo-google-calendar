package app

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/earlgray283/gakujo-google-calendar/assets"
	"github.com/earlgray283/gakujo-google-calendar/gakujo"
	"github.com/skratchdot/open-golang/open"
)

type AuthFormInfo struct {
	Username  string
	Password  string
	Logincode string
}

type templates struct {
	URL       string
	ErrorHTML template.HTML
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
		if err := t.Execute(rw, templates{URL: url}); err != nil {
			http.Error(rw, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/regist", func(rw http.ResponseWriter, r *http.Request) {
		m1, _ := regexp.MatchString("^[a-zA-Z0-9]+$", r.FormValue("username"))
		m2, _ := regexp.MatchString("^[a-zA-Z0-9]+$", r.FormValue("password"))
		if !m1 || !m2 {
			rw.WriteHeader(http.StatusBadRequest)
			_ = t.Execute(rw, templates{
				ErrorHTML: template.HTML(`<div class="alert alert-danger" role="alert">静大IDか静大パスワードの形式が正しくありません。</div>`),
				URL:       url,
			})
			return
		}

		c := gakujo.NewClient()
		if err := c.Login(r.FormValue("username"), r.FormValue("password")); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			_ = t.Execute(rw, templates{
				ErrorHTML: template.HTML(`<div class="alert alert-danger" role="alert">静大IDか静大パスワードが間違っています。</div>`),
				URL:       url,
			})
			return
		}

		authFormInfo = AuthFormInfo{
			Username:  r.FormValue("username"),
			Password:  r.FormValue("password"),
			Logincode: r.FormValue("logincode"),
		}

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
