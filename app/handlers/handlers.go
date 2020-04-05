package handlers

import (
  "bytes"
  "strconv"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
  "net/url"
  "html/template"
	"os"
	"time"
  "log"

	"github.com/RollMan/ac2manager/app/db"
	"github.com/RollMan/ac2manager/app/models"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var allColumnsOfEvent = "id, startdate, track, weatherRandomness, P_hourOfDay, P_timeMultiplier, P_sessionDurationMinute, Q_hourOfDay, Q_timeMultiplier, Q_sessionDurationMinute, R_hourOfDay, R_timeMultiplier, R_sessionDurationMinute, pitWindowLengthSec, isRefuellingAllowedInRace, mandatoryPitstopCount, isMandatoryPitstopRefuellingRequired, isMandatoryPitstopTyreChangeRequired, isMandatoryPitstopSwapDriverRequired, tyreSetCount"

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
	row := db.Db.QueryRow(fmt.Sprintf("SELECT %s FROM events WHERE events.startdate >= CONVERT(?, DATETIME) ORDER BY startdate ASC;", allColumnsOfEvent), now)
  err := row.Scan(&event.Id, &event.Startdate, &event.Track, &event.WeatherRandomness, &event.P_hourOfDay, &event.P_timeMultiplier, &event.P_sessionDurationMinute, &event.Q_hourOfDay, &event.Q_timeMultiplier, &event.Q_sessionDurationMinute, &event.R_hourOfDay, &event.R_timeMultiplier, &event.R_sessionDurationMinute, &event.PitWindowLengthSec, &event.IsRefuellingAllowedInRace, &event.MandatoryPitstopCount, &event.IsMandatoryPitstopRefuellingRequired, &event.IsMandatoryPitstopTyreChangeRequired, &event.IsMandatoryPitstopSwapDriverRequired, &event.TyreSetCount)


  jst := time.FixedZone("JST", 9*60*60)
  event.Startdate = event.Startdate.In(jst)

  var isNextRace bool = true
  if err != nil {
    if err == sql.ErrNoRows {
      isNextRace = false
    }else{
      returnInternalServerError(w, err)
      return
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
  var writeBuf bytes.Buffer
  t := template.Must(template.ParseFiles("./template/login.html"))
  err := t.Execute(&writeBuf, nil)
  if err != nil {
    returnInternalServerError(w, err)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(writeBuf.Bytes())
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	var err error
  r.ParseForm()

	userid := r.Form.Get("userid")
	pw := r.Form.Get("pw")

	var user models.User
	user, err = checkUserPw([]byte(userid), []byte(pw))
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

func AdminHandler(w http.ResponseWriter, r *http.Request, token *TokenClaims) {
  var events []models.Event
	rows, err := db.Db.Query(fmt.Sprintf("SELECT %s FROM events ORDER BY startdate DESC;", allColumnsOfEvent))
  if err != nil {
    returnInternalServerError(w, err)
    return
  }


  for rows.Next() {
    var event models.Event
    err := rows.Scan(&event.Id, &event.Startdate, &event.Track, &event.WeatherRandomness, &event.P_hourOfDay, &event.P_timeMultiplier, &event.P_sessionDurationMinute, &event.Q_hourOfDay, &event.Q_timeMultiplier, &event.Q_sessionDurationMinute, &event.R_hourOfDay, &event.R_timeMultiplier, &event.R_sessionDurationMinute, &event.PitWindowLengthSec, &event.IsRefuellingAllowedInRace, &event.MandatoryPitstopCount, &event.IsMandatoryPitstopRefuellingRequired, &event.IsMandatoryPitstopTyreChangeRequired, &event.IsMandatoryPitstopSwapDriverRequired, &event.TyreSetCount)
    if err != nil {
      returnInternalServerError(w, err)
      return
    }
    jst := time.FixedZone("JST", 9*60*60)
    event.Startdate = event.Startdate.In(jst)
    events = append(events, event)
  }

  var writeBuf bytes.Buffer
  t := template.Must(template.ParseFiles("./template/admin.html"))
  data := struct {
    EventTableRows []models.Event
  }{events}
  err = t.Execute(&writeBuf, data)
  if err != nil {
    returnInternalServerError(w, err)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(writeBuf.Bytes())
}

func AuthMiddleware(next HttpHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("In authmiddle handlerfunc")
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

func decodeAddForm(form url.Values) (models.Event, error) {
  for k, v := range form{
    if len(v) == 0 || v[0] == ""{
      return models.Event{}, fmt.Errorf("No field for key %v.", k)
    }
  }
  startdate_day := form.Get("startdate_day")
  startdate_time := form.Get("startdate_time")
  startdate, err := time.Parse("2006-01-02 15:04", startdate_day + " " + startdate_time)
  if err != nil {
    return models.Event{}, err
  }

  atoi := func (s string) int {
    r, err := strconv.Atoi(s)
    if err != nil {
      log.Panic(fmt.Sprintf("PANIC: Failed to parse form of %s into integer: %v", s, err.Error()))
    }
    return r
  }
  parseBool := func (s string) bool {
    r, err := strconv.ParseBool(s)
    if err != nil {
      log.Panic(fmt.Sprintf("PANIC: Failed to parse form of %s into bool: %v", s, err.Error()))
    }
    return r
  }
  return models.Event {
    Startdate: startdate,
    Track: form.Get("track"),
    WeatherRandomness: atoi(form.Get("weatherRandomness")),
                P_hourOfDay:                            atoi(form.Get("P_hourOfDay")),
                P_timeMultiplier:                       atoi(form.Get("P_timeMultiplier")),
                P_sessionDurationMinute:                atoi(form.Get("P_sessionDurationMinute")),
                Q_hourOfDay:                            atoi(form.Get("Q_hourOfDay")),
                Q_timeMultiplier:                       atoi(form.Get("Q_timeMultiplier")),
                Q_sessionDurationMinute:                atoi(form.Get("Q_sessionDurationMinute")),
                R_hourOfDay:                            atoi(form.Get("R_hourOfDay")),
                R_timeMultiplier:                       atoi(form.Get("R_timeMultiplier")),
                R_sessionDurationMinute:                atoi(form.Get("R_sessionDurationMinute")),
                PitWindowLengthSec:                     atoi(form.Get("pitWindowLengthSec")),
                IsRefuellingAllowedInRace:              parseBool(form.Get("isRefuellingAllowedInRace")),
                MandatoryPitstopCount:                  atoi(form.Get("mandatoryPitstopCount")),
                IsMandatoryPitstopRefuellingRequired:   parseBool(form.Get("isMandatoryPitstopRefuellingRequired")),
                IsMandatoryPitstopTyreChangeRequired:   parseBool(form.Get("isMandatoryPitstopTyreChangeRequired")),
                IsMandatoryPitstopSwapDriverRequired:   parseBool(form.Get("isMandatoryPitstopSwapDriverRequired")),
                TyreSetCount:                           atoi(form.Get("tyreSetCount")),
  }, nil
}

func AddHandler(w http.ResponseWriter, r *http.Request, token *TokenClaims) {
  r.ParseForm()
  event, err := decodeAddForm(r.Form)
  if err != nil {
    returnInternalServerError(w, err)
    return
  }

  log.Printf("%v", event)

  rows, err := db.Db.Query("SELECT * FROM events")
  if err != nil {
    returnInternalServerError(w, err)
    return
  }

  var isDupicate = false
  var duplicating models.Event
  for rows.Next() {
    var target models.Event
    err := rows.Scan(&target.Id, &target.Startdate, &target.Track, &target.WeatherRandomness, &target.P_hourOfDay, &target.P_timeMultiplier, &target.P_sessionDurationMinute, &target.Q_hourOfDay, &target.Q_timeMultiplier, &target.Q_sessionDurationMinute, &target.R_hourOfDay, &target.R_timeMultiplier, &target.R_sessionDurationMinute, &target.PitWindowLengthSec, &target.IsRefuellingAllowedInRace, &target.MandatoryPitstopCount, &target.IsMandatoryPitstopRefuellingRequired, &target.IsMandatoryPitstopTyreChangeRequired, &target.IsMandatoryPitstopSwapDriverRequired, &target.TyreSetCount)
    if err != nil {
      returnInternalServerError(w, err)
      return
    }
    adding_start := event.Startdate
    adding_end := adding_start.Add(time.Minute * time.Duration(event.P_sessionDurationMinute + event.Q_sessionDurationMinute + event.R_sessionDurationMinute + 5))
    target_start := target.Startdate
    target_end := target_start.Add(time.Minute * time.Duration(target.P_sessionDurationMinute + target.Q_sessionDurationMinute + target.R_sessionDurationMinute + 5))
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

  ins, err := db.Db.Prepare(fmt.Sprintf("INSERT INTO events (startdate, track, weatherRandomness, P_hourOfDay, P_timeMultiplier, P_sessionDurationMinute, Q_hourOfDay, Q_timeMultiplier, Q_sessionDurationMinute, R_hourOfDay, R_timeMultiplier, R_sessionDurationMinute, pitWindowLengthSec, isRefuellingAllowedInRace, mandatoryPitstopCount, isMandatoryPitstopRefuellingRequired, isMandatoryPitstopTyreChangeRequired, isMandatoryPitstopSwapDriverRequired, tyreSetCount) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"))
  if err != nil {
    returnInternalServerError(w, err)
    return
  }
  ins.Exec(event.Startdate, event.Track, event.WeatherRandomness, event.P_hourOfDay, event.P_timeMultiplier, event.P_sessionDurationMinute, event.Q_hourOfDay, event.Q_timeMultiplier, event.Q_sessionDurationMinute, event.R_hourOfDay, event.R_timeMultiplier, event.R_sessionDurationMinute, event.PitWindowLengthSec, event.IsRefuellingAllowedInRace, event.MandatoryPitstopCount, event.IsMandatoryPitstopRefuellingRequired, event.IsMandatoryPitstopTyreChangeRequired, event.IsMandatoryPitstopSwapDriverRequired, event.TyreSetCount)

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(fmt.Sprintf("OK %v", event)))
  log.Print(event)
}

func isNoDuplicate(a_start, a_end, b_start, b_end time.Time) bool {
  return a_end.Before(b_start) || b_end.Before(a_start)
}

func AddEventHandler(w http.ResponseWriter, r *http.Request, token *TokenClaims){
  var writeBuf bytes.Buffer
  t := template.Must(template.ParseFiles("./template/add.html"))
  data := map[string]string{}
  err := t.Execute(&writeBuf, data)
  if err != nil {
    returnInternalServerError(w, err)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(writeBuf.Bytes())
}
