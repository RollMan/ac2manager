package apiHandlers

import (
  "os"
  "log"
  "database/sql"
  "fmt"
  "errors"
  "time"
  "encoding/json"
	"net/http"
  "io/ioutil"
	"github.com/RollMan/ac2manager/app/db"
	"github.com/RollMan/ac2manager/app/models"
  jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	Attribute int
	jwt.StandardClaims
}

func ServerStatusHandler(w http.ResponseWriter, r *http.Request){
}

func RacesHandler(w http.ResponseWriter, r *http.Request){
  var events []models.Event
  {
    _, err := db.DbMap.Select(&events, "SELECT * FROM events ORDER BY startdate DESC;")

    if err != nil {
      log.Printf("%v\n", err)
      body := fmt.Sprintf("%v\n", err)
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(body))
      return
    }
  }

  {
    body, err := json.Marshal(events)
    if err != nil {
      body := []byte(fmt.Sprintf("%v\n", err))
      log.Println(body)
      w.WriteHeader(http.StatusInternalServerError)
      w.Write(body)
      return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(body)
  }
}

func UpcomingRaceHandler(w http.ResponseWriter, r *http.Request){
  event := make([]models.Event, 1)
  var isNextRace bool = true
  now := time.Now()
  err := db.DbMap.SelectOne(&event[0], "SELECT * FROM events WHERE events.startdate >= CONVERT(?, DATETIME) ORDER BY startdate ASC;", now)

  if err != nil {
    if err == sql.ErrNoRows {
      isNextRace = false
    }else{
      w.WriteHeader(http.StatusInternalServerError)
      body := fmt.Sprintf("%v\n", err)
      w.Write([]byte(body))
      return
    }
  }

  if !isNextRace {
    emptyEvent := make([]models.Event, 0)
    body, err := json.Marshal(emptyEvent)
    if err != nil {
      body := []byte(fmt.Sprintf("%v\n", err))
      log.Println(body)
      w.WriteHeader(http.StatusInternalServerError)
      w.Write(body)
      return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(body))
  }else{
    body, err := json.Marshal(event)
    if err != nil {
      body := []byte(fmt.Sprintf("%v\n", err))
      log.Println(body)
      w.WriteHeader(http.StatusInternalServerError)
      w.Write(body)
      return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(body))
  }
}

func PastRacesHandler(w http.ResponseWriter, r *http.Request){
  var events []models.Event
  now := time.Now()
  _, err := db.DbMap.Select(&events, "SELECT * FROM events WHERE events.startdate <= CONVERT(?, DATETIME ORDER BY startdate ASC;", now)
  if err == nil {
    if err == sql.ErrNoRows {
      events = make([]models.Event, 0)
    }else{
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(err.Error()))
      return
    }
  }

  body, err := json.Marshal(events)
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
  }
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(body))
}

func FutureRacesHandler(w http.ResponseWriter, r *http.Request){
  var events []models.Event
  now := time.Now()
  _, err := db.DbMap.Select(&events, "SELECT * FROM events WHERE events.startdate > CONVERT(?, DATETIME ORDER BY startdate ASC;", now)
  if err == nil {
    if err == sql.ErrNoRows {
      events = make([]models.Event, 0)
    }else{
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(err.Error()))
      return
    }
  }

  body, err := json.Marshal(events)
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
  }
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(body))
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
  var err error
  if r.Header.Get("Content-Type") != "application/json" {
  }else{
    log.Println("Invalid, not application/json request for login received.")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  buf, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  var req_user models.User
  err = json.Unmarshal(buf, &req_user)
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  userid := req_user.UserID
  pwhash := req_user.PWHash

  err = checkUserPw([]byte(userid), []byte(pwhash))
	if err != nil {
		switch err.(type) {
		case *models.NoSuchUserError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Such user does not exist."))
      log.Printf("ERROR: Login request submitted of non-existing user %s. %v", userid, err.Error())
		case *models.NoMatchingPasswordError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong password."))
      log.Printf("ERROR: Login request submitted wrong password for user %s. %v", userid, err.Error())
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unknown error. Contact administrator."))
      log.Printf("ERROR: Unknown error when checking password. %v", err.Error())
		}
		return
	}

  now := time.Now()
	claims := &TokenClaims{
		req_user.Attribute,
		jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * 6).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Audience:  string(userid),
		},
	}

	signingKey := []byte(os.Getenv("JWT_SIGNING_KEY")) // Specify in `.env`
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
    log.Printf("ERROR: Error occured when signing JWT. %v", err.Error())
		return
	}

  res := map[string]string {"jwt": ss}
  body, err := json.Marshal(res)
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func checkUserPw(userid []byte, pw []byte) error {
	var user models.User
	var err error
  err = db.DbMap.SelectOne(&user, "SELECT * FROM users WHERE userid=?;", userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.NoSuchUserError{}
		} else {
      return errors.New(fmt.Sprintf("Unknown Error: %v", err))
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PWHash), pw); err != nil {
		return &models.NoMatchingPasswordError{}
	}
	return nil
}
