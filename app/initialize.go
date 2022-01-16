package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
)
func main() {

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
		username := r.FormValue("username")
		password := r.FormValue("password")
		token := r.FormValue("token")
		if len(username) == 0 || len(password) == 0 || len(token) == 0 {
    		http.Error(rw, "username, password, token must not be empty", http.StatusBadRequest)
    	return
		}
		fmt.Println(username, password, token)
	})

	log.Println("Listening on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}