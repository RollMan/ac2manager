package handlers

import (
  "bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
  "html/template"
	"os"
	"time"
  "log"

	"github.com/RollMan554/ac2manager/app/db"
	"github.com/RollMan554/ac2manager/app/models"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type HttpHandler func(http.ResponseWriter, *http.Request, *TokenClaims)
type TokenClaims struct {
	Attribute int
	jwt.StandardClaims
}

func returnInternalServerError(w http.ResponseWriter, err error){
    w.WriteHeader(http.StatusInternalServerError)
    t_internalServerError := template.Must(template.ParseFiles("./template/StatusInternalServerError.html"))
    data := map[string]string{
      "Message": err.Error(),
    }
    t_internalServerError.Execute(w, data)
    log.Print(err.Error())
    return
}

func AboutHandler(w http.ResponseWriter, r *http.Request){
  t := template.Must(template.ParseFiles("./template/about.html"))
  data := map[string]string{"a":"a"}
  var writeBuf bytes.Buffer
  err := t.Execute(&writeBuf, data)
  if err != nil {
    returnInternalServerError(w, err)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(writeBuf.Bytes())
}


func RootHandler(w http.ResponseWriter, r *http.Request) {
  // Find next event
  // SELECT * FROM events WHERE events.startdate >= CONVERT('1999-01-00', DATETIME) ORDER BY startdate DESC;
  var event models.Event
  now := time.Now()
  columns := "id, startdate, track, weatherRandomness, P_hourOfDay, P_timeMultiplier, P_sessionDurationMinute, Q_hourOfDay, Q_timeMultiplier, Q_sessionDurationMinute, R_hourOfDay, R_timeMultiplier, R_sessionDurationMinute, pitWindowLengthSec, isRefuellingAllowedInRace, mandatoryPitstopCount, isMandatoryPitstopRefuellingRequired, isMandatoryPitstopTyreChangeRequired, isMandatoryPitstopSwapDriverRequired, tyreSetCount"
	row := db.Db.QueryRow(fmt.Sprintf("SELECT %s FROM events WHERE events.startdate >= CONVERT(?, DATETIME) ORDER BY startdate ASC;", columns), now)
  err := row.Scan(&event.Id, &event.Startdate, &event.Track, &event.WeatherRandomness, &event.P_hourOfDay, &event.P_timeMultiplier, &event.P_sessionDurationMinute, &event.Q_hourOfDay, &event.Q_timeMultiplier, &event.Q_sessionDurationMinute, &event.R_hourOfDay, &event.R_timeMultiplier, &event.R_sessionDurationMinute, &event.PitWindowLengthSec, &event.IsRefuellingAllowedInRace, &event.MandatoryPitstopCount, &event.IsMandatoryPitstopRefuellingRequired, &event.IsMandatoryPitstopTyreChangeRequired, &event.IsMandatoryPitstopSwapDriverRequired, &event.TyreSetCount)


  jst := time.FixedZone("JST", 9*60*60)
  event.Startdate = event.Startdate.In(jst)

  var isNextRace bool = true
  if err != nil {
    if err == sql.ErrNoRows {
      isNextRace = false
    }else{
      returnInternalServerError(w, err)
    }
  }

  var writeBuf bytes.Buffer
  if isNextRace {
    data := models.NextRaceData{
      event,
      "SERVER STATUS ICON",
      "SERVER STATUS STATEMENT",
    }
    t := template.Must(template.ParseFiles("./template/index.html", "./template/upcoming_race_configure.html"))
    err = t.Execute(&writeBuf, data)
    if err != nil {
      returnInternalServerError(w, err)
      return
    }
  }else{
    data := map[string]string{}
    t := template.Must(template.ParseFiles("./template/index.html", "./template/no_upcoming_race.html"))
    err = t.Execute(&writeBuf, data)
    if err != nil {
      returnInternalServerError(w, err)
      return
    }
  }

  w.WriteHeader(http.StatusOK)
  w.Write(writeBuf.Bytes())
}

func LoginGetHandler(w http.ResponseWriter, r *http.Request) {

}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var req_json models.Login
	if err := json.NewDecoder(r.Body).Decode(&req_json); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unknown error. Couldn't decode JSON.\n%s\n", err)))
    log.Printf("ERROR: Couldn't parse JSON in LoginPostHandler. %v", err.Error())
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
      log.Printf("ERROR: Login request submitted non-existing user %s. %v", userid, err.Error())
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
    log.Printf("ERROR: Error occured when signing JWT. %v", err.Error())
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

func AdminHandler(w http.ResponseWriter, r *http.Request, t *TokenClaims) {
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

    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte("Internal Server Error"))
      log.Printf("ERROR: Failed to parse JWT. %v", err.Error())
      return
    }

		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			next(w, r, claims)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
      log.Printf("ERROR: Token was not vaild.")
		}
	})
}

func checkUserPw(userid []byte, pw []byte) (models.User, error) {
	var user models.User
	var err error
	row := db.Db.QueryRow("SELECT * FROM users WHERE userid=?;", userid)
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
