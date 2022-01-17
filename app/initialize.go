package app

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"
)

type User struct {
	Username  string
	Password  string
	Logincode string
}

/*func main(){
	UserInfo := GetUserInfoFromBrowser()
	fmt.Println(UserInfo.Username, UserInfo.Password, UserInfo.Token)
}*/

func GetUserInfoFromBrowser(url string) User {
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

	mux.HandleFunc("/regist", func(rw http.ResponseWriter, r *http.Request) {
		UserInfo.Username = r.FormValue("username")
		UserInfo.Password = r.FormValue("password")
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
		log.Fatalln("Server closed with error:", err)
	}

	return UserInfo

}
