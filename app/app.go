package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
  "errors"
  "encoding/json"

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
  var req_json models.Login
  if err := json.NewDecoder(r.Body).Decode(&req_json); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(fmt.Sprintf("Unknown error. Couldn't decode JSON.\n%s\n", err)))
    return
  }

  userid := req_json.UserID
  pw := req_json.Password

  err = checkUserPw([]byte(userid), []byte(pw))
  if err != nil {
    switch err.(type) {
    case *models.NoSuchUserError:
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Such user does not exist."))
    case *models.NoMatchingPasswordError:
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

func checkUserPw(userid []byte, pw []byte) (error){
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

  if err := bcrypt.CompareHashAndPassword([]byte(user.PWHash), pw); err != nil {
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

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	// r.HandleFunc("/admin", AuthMiddle)
  
	log.Fatal(http.ListenAndServe(":8000", r))
}
