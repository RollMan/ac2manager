package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
  "errors"

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
  var err error

	r.ParseForm()
	userid := r.PostForm.Get("userid")
	pw := r.PostForm.Get("pw")

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		return
	}

  err = checkUserPw(userid, string(hash))
  if err != nil {
    switch err.(type) {
    case models.NoSuchUserError:
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Such user does not exist."))
    case models.NoMatchingPasswordError:
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Wrong password."))
    default:
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte("Unknown error. Contact administrator."))
    }
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("WIP"))
}

func checkUserPw(userid string, pwhash string) (error){
  var user models.User
  var err error
  row := db.QueryRow("SELECT * FROM users WHERE userid=?;", userid)
  err = row.Scan(&user.UserID, &user.PWHash, &user.Attribute)

  if err != nil {
    if err == sql.ErrNoRows {
      return &models.NoSuchUserError{}
    } else {
      return errors.New("Unknown Error")
    }
  }

  if user.PWHash != pwhash {
    return &models.NoMatchingPasswordError{}
  }
  return nil
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
	}

	fmt.Printf("Hello World\n")
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	// r.HandleFunc("/admin", AuthMiddle)
  
	// log.Fatal(http.ListenAndServe(":8000", r))
}
