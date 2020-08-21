package jobmng

import (
  "time"
  "log"
  "github.com/RollMan/ac2manager/ec2ctl/db"
  "github.com/RollMan/ac2manager/app/models"
  "github.com/RollMan/ac2manager/ec2ctl/confjson"
)

var queue []models.Event

func InitQueue(){
  queue = make([]models.Event, 0)
}

func FindJobs(t time.Time) {
  targetInMinute := t.Truncate(time.Minute)
  events := selectJobsByDate(targetInMinute)
  for _, e := range events {
    queue = append(queue, e)
  }
}

func selectJobsByDate(t time.Time){
  var events []models.Event
  t1 := t
  t2 := t.Add(time.Minute)
  _, err := db.DbMap.Select(&events, "SELECT * FROM events WHERE CONVERT(?, DATETIME) <= events.startdate and events.startdate < CONVERT(?, DATETIME)", t1, t2)

  if err != nil {
    log.Fatal(err)
  }
  return events
}

func RunQueue(){
  // TODO
  // create configure json
  // deploy server
}
