package apiHandlers

import (
	// "github.com/RollMan/ac2manager/app/db"
	"github.com/RollMan/ac2manager/app/models"
	guuid "github.com/google/uuid"
	"time"
)

func updateEc2ctl(e models.Event) error {
	start := models.StartStop{
		Id:        models.UUID(guuid.New()),
    EventID: e.Id
		Op: "Start",
		Datetime:  e.Startdate,
	}

  stoptime := calcStopTime(e)

  stop := models.StartStop {
    Id: models.UUID(guuid.New()),
    EventID: e.Id
    Op: "Stop",
    Datetime: stoptime,
  }

  db.DbMap.Insert(&start)
  db.DbMap.Insert(&stop)
  return nil
}

func calcStopTime(e models.Event) time.Time {
  return e.Startdate.Add(time.Duration(e.P_sessionDurationMinute) * time.Minute + time.Duration(e.Q_sessionDurationMinute) * time.Minute + time.Duration(e.R_sessionDurationMinute) + time.Minute + 30 * time.Minute)

}
