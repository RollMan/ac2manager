package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

  "github.com/RollMan554/ac2manager/app/models"
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

func checkUserPw(id []byte, pw []byte){
  var user models.User
  var err error
  row := db.QueryRow("SELECT * FROM users WHERE userid=?;", id)
  err = row.Scan(&user.UserID, &user.PWHash, &user.Attribute)

  if err != nil {
    if err == sql.ErrNoRows {
      fmt.Printf("No such user: %s\n", id)
    } else {
      log.Fatal(err)
    }
  }
  fmt.Printf("PW: %s %d\n", user.PWHash, user.Attribute)
}

func main() {
	{
		dsn := fmt.Sprintf("%s:%s@/ac2?charset=utf8", os.Getenv("AC2_DB_USERNAME"), os.Getenv("AC2_DB_PASSWORD"))
		var err error
		db, err = sql.Open("mysql", dsn)
    if err != nil {
      fmt.Print("DB Opening Error")
      log.Fatal(err)
    }
    err = db.Ping()
    if err != nil {
      fmt.Print("DB Ping Error")
      log.Fatal(err)
    }
    fmt.Print("DB OK")
    checkUserPw([]byte("admin"), []byte("password"))
    os.Exit(0)
	}

	fmt.Printf("Hello World\n")
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	// r.HandleFunc("/admin", AuthMiddle)
  
	// log.Fatal(http.ListenAndServe(":8000", r))
}
