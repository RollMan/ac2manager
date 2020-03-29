package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RollMan554/ac2manager/app/models"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type HttpHandler func(http.ResponseWriter, *http.Request, *TokenClaims)
type TokenClaims struct {
	Attribute int
	jwt.StandardClaims
}

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

	var user models.User
	user, err = checkUserPw([]byte(userid), []byte(pw))
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

	now := time.Now()
	claims := &TokenClaims{
		user.Attribute,
		jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * 6).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Audience:  userid,
		},
	}

	signingKey := []byte(os.Getenv("JWT_SIGNING_KEY")) // Specify in `.env`
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	cookie := &http.Cookie{
		Name:  "jwt",
		Value: ss,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success!"))
}

func checkUserPw(userid []byte, pw []byte) (models.User, error) {
	var user models.User
	var err error
	row := db.QueryRow("SELECT * FROM users WHERE userid=?;", userid)
	err = row.Scan(&user.UserID, &user.PWHash, &user.Attribute)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, &models.NoSuchUserError{}
		} else {
			return models.User{}, errors.New("Unknown Error")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PWHash), pw); err != nil {
		return models.User{}, &models.NoMatchingPasswordError{}
	}
	return user, nil
}

func adminHandler(w http.ResponseWriter, r *http.Request, t *TokenClaims) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Secret admin page"))
}

func AuthMiddleware(next HttpHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("jwt")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("no token"))
			return
		}

		tokenstring := tokenCookie.Value

		token, err := jwt.ParseWithClaims(tokenstring, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
		})

		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			next(w, r, claims)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})
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
	r.HandleFunc("/admin", AuthMiddleware(adminHandler))

	log.Fatal(http.ListenAndServe(":8000", r))
}
