package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gorilla!\n"))
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {

}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	pw := r.PostForm.Get("pw")

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		return
	}

	log.Printf("%s %s", id, hash)
}

func main() {
	fmt.Printf("Hello World\n")
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	// r.HandleFunc("/admin", AuthMiddle)
	log.Fatal(http.ListenAndServe(":8000", r))
}
