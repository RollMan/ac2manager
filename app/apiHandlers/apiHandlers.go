package apiHandlers

import (
  "strings"
  "reflect"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RollMan/ac2manager/app/db"
	"github.com/RollMan/ac2manager/app/models"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
  "github.com/mholt/binding"
)

type TokenClaims struct {
	Attribute int
	jwt.StandardClaims
}

func ParseJSONBody(r *http.Request, res interface{}) error {
	var err error

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, res)
	if err != nil {
		return err
	}
	return nil

}

func ServerStatusHandler(w http.ResponseWriter, r *http.Request) {
}

func RacesHandler(w http.ResponseWriter, r *http.Request) {
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

func UpcomingRaceHandler(w http.ResponseWriter, r *http.Request) {
	event := make([]models.Event, 1)
	var isNextRace bool = true
	now := time.Now()
	err := db.DbMap.SelectOne(&event[0], "SELECT * FROM events WHERE events.startdate >= CONVERT(?, DATETIME) ORDER BY startdate ASC;", now)

	if err != nil {
		if err == sql.ErrNoRows {
			isNextRace = false
		} else {
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
	} else {
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

func PastRacesHandler(w http.ResponseWriter, r *http.Request) {
	var events []models.Event
	now := time.Now()
	_, err := db.DbMap.Select(&events, "SELECT * FROM events WHERE events.startdate <= CONVERT(?, DATETIME) ORDER BY startdate ASC;", now)
	if err != nil {
		if err == sql.ErrNoRows {
			events = make([]models.Event, 0)
		} else {
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

func FutureRacesHandler(w http.ResponseWriter, r *http.Request) {
	var events []models.Event
	now := time.Now()
	_, err := db.DbMap.Select(&events, "SELECT * FROM events WHERE events.startdate > CONVERT(?, DATETIME) ORDER BY startdate ASC;", now)
	if err != nil {
		if err == sql.ErrNoRows {
			events = make([]models.Event, 0)
		} else {
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("Invalid, not application/json request for login received.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req_user models.User
	err = ParseJSONBody(r, &req_user)
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

	res := map[string]string{"jwt": ss}
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

func AddRaceHandler(w http.ResponseWriter, r *http.Request, token *models.TokenClaims) {
	var err error
	var event models.Event
  if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
    err := binding.Bind(r, &event)
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    log.Println(r)
    log.Println(event)
  }else if r.Header.Get("Content-Type") == "application/json" {
    err = ParseJSONBody(r, &event)
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
  }else {
    log.Println("Invalid, not application/json request for login received.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var events []models.Event
	_, err = db.DbMap.Select(&events, "SELECT * FROM events")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var isDupicate = false
	var duplicating models.Event
	for _, target := range events {
		adding_start := event.Startdate
		adding_end := adding_start.Add(time.Minute * time.Duration(event.P_sessionDurationMinute+event.Q_sessionDurationMinute+event.R_sessionDurationMinute+5))
		target_start := target.Startdate
		target_end := target_start.Add(time.Minute * time.Duration(target.P_sessionDurationMinute+target.Q_sessionDurationMinute+target.R_sessionDurationMinute+5))
		if isNoDuplicate(adding_start, adding_end, target_start, target_end) == false {
			isDupicate = true
			break
		}
	}
	if isDupicate {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("The adding event duration duplicates with the other event.\nadding: %v\ntarget: %v", event, duplicating)))
		return
	}
	err = db.DbMap.Insert(&event)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res map[string]interface{} = map[string]interface{}{"success": true}

	body, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func isNoDuplicate(a_start, a_end, b_start, b_end time.Time) bool {
	return a_end.Before(b_start) || b_end.Before(a_start)
}

func RemoveRaceHandler(w http.ResponseWriter, r *http.Request, token *models.TokenClaims) {
	var err error
	var event models.Event
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("Invalid, not application/json request for login received.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = ParseJSONBody(r, &event)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	count, err := db.DbMap.Delete(&event)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := map[string]interface{}{"count": count}
	body, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func SchemaHandler(w http.ResponseWriter, r *http.Request){
  var event models.Event
  t := reflect.TypeOf(event)

  schema := make(map[string]string) // [keyname]type

  for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    json_keyname := field.Tag.Get("json")
    typename := field.Type.Name()
    schema[json_keyname] = typename
  }

  body,err := json.Marshal(schema)
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(body))
}
