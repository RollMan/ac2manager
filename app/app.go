package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

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
	{
		dsn := fmt.Sprintf("%s:%s@ac2", os.Getenv("AC2_DB_USERNAME"), os.Getenv("AC2_DB_PASSWORD"))
		var err error
		db, err := sql.Open("mysql", dsn)
		log.Fatal(db)
		log.Fatal(err)
	}
	fmt.Printf("Hello World\n")
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	// r.HandleFunc("/admin", AuthMiddle)
	log.Fatal(http.ListenAndServe(":8000", r))
}
