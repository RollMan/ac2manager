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

func TestFindJobs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error %s occuered when opening db mock", err)
	}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbMap.AddTableWithName(models.Event{}, "events").SetKeys(true, "id")
	defer dbMap.Db.Close()

	emptyRows := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"id", "startdate"})
	}

	time1 := time.Date(2020, 5, 3, 23, 0, 0, 0, time.UTC)
	time2 := time.Date(2020, 5, 5, 11, 30, 0, 0, time.UTC)

	cases := []struct {
		description string
		rows        *sqlmock.Rows
		time        time.Time
		expected    []JobQueue
	}{
		{
			description: "Test when current time matches an event",
			rows:        emptyRows().AddRow(3, time1),
			time:        time1,
			expected: []JobQueue{
				{
					JobType: Start,
					Event:   models.Event{Id: 3, Startdate: time1},
				},
				{
					JobType: Stop,
					Event:   models.Event{Id: 3, Startdate: time1.Add(time.Minute * 10)}, // Assume that non-initialized time.Time is 0.
				},
			},
		},
		{
			description: "Test when current time matches plural events",
			rows:        emptyRows().AddRow(0, time1).AddRow(5, time1),
			time:        time1,
			expected: []JobQueue{
				{
					JobType: Start,
					Event:   models.Event{Id: 0, Startdate: time1},
				},
				{
					JobType: Stop,
					Event:   models.Event{Id: 0, Startdate: time1.Add(time.Minute * 10)},
				},
				{
					JobType: Start, Event: models.Event{Id: 5, Startdate: time1},
				},
				{
					JobType: Stop, Event: models.Event{Id: 5, Startdate: time1.Add(time.Minute * 10)},
				},
			},
		},
		{
			description: "Test when current time does not maches any events",
			rows:        emptyRows(),
			time:        time2,
			expected:    []JobQueue{},
		},
	}

	for _, c := range cases {
		mock.ExpectQuery(`SELECT \* FROM events`).
			WithArgs(c.time, c.time.Add(time.Minute)).
			WillReturnRows(c.rows)

		res := FindJobs(c.time, []JobQueue{}, dbMap)

		if len(res) != len(c.expected) {
			t.Errorf("Error in a test.\ntestcase description: %s\nThe number of result %d must %d but not.\nres:%v\nexpected:%v\n", c.description, len(res), len(c.expected), res, c.expected)
		}

		for i, _ := range res {
			if res[i] != c.expected[i] {
				t.Errorf("Error in a test.\ntestcase description: %s\nThe result unmatch.\nres:%v\nexpected:%v\n", c.description, res[i], c.expected[i])
			}
		}
	}
}
