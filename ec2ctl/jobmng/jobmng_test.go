package jobmng

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RollMan/ac2manager/app/models"
	"github.com/go-gorp/gorp"
	"testing"
	"time"
)

func TestSelectJobsByDate(t *testing.T) {
	// queue := InitQueue()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error %s occuered when opening db mock", err)
	}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbMap.AddTableWithName(models.Event{}, "events").SetKeys(true, "id")
	defer dbMap.Db.Close()

	target_time := time.Date(2020, 9, 12, 10, 30, 0, 0, time.UTC)
	target_Time2 := target_time.Add(time.Minute)
	{
		row := sqlmock.NewRows([]string{"id", "startdate"}).AddRow(0, target_time)
		expected := models.Event{Id: 0, Startdate: target_time}
		mock.ExpectQuery(`SELECT \* FROM events`).
			WithArgs(target_time, target_Time2).
			WillReturnRows(row)

		events := selectJobsByDate(target_time, dbMap)
		if events[0] != expected {
			t.Errorf("invalid result")
		}
	}

	{
		mock.ExpectQuery(`SELECT \* FROM events`).
			WithArgs(target_time, target_Time2).
			WillReturnRows(sqlmock.NewRows([]string{"id", "startdate"}))

		events := selectJobsByDate(target_time, dbMap)
		if len(events) != 0 {
			t.Errorf("invalid result")
		}
	}
}
