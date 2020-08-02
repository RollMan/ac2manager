package apiHandlers

import (
  "log"
  "database/sql"
  "fmt"
  "time"
  "encoding/json"
	"net/http"
	"github.com/RollMan/ac2manager/app/db"
	"github.com/RollMan/ac2manager/app/models"
)

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
