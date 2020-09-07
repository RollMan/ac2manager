package jobmng

import (
	"github.com/RollMan/ac2manager/app/models"
	"github.com/RollMan/ac2manager/ec2ctl/confjson"
	"github.com/RollMan/ac2manager/ec2ctl/db"
	"log"
	"time"
)

type JobType int

const (
	_ JobType = iota
	Start
	Stop
)

type jobQueue struct {
	JobType        JobType
	Event          models.Event
	LaunchSchedule time.Time
}

var queue []jobQueue

func InitQueue() {
	queue = make([]jobQueue, 0)
}

func FindJobs(t time.Time) {
	targetInMinute := t.Truncate(time.Minute)
	events := selectJobsByDate(targetInMinute)
	for _, e := range events {
		queue = append(queue, jobQueue{Start, e, e.Startdate})
		extra := time.Minute * 10
		enddate := e.Startdate.Add(time.Minute*time.Duration(e.P_sessionDurationMinute+e.Q_sessionDurationMinute+e.R_sessionDurationMinute) + extra)
		queue = append(queue, jobQueue{Stop, e, enddate})
	}
}

func selectJobsByDate(t time.Time) []models.Event {
	var events []models.Event
	t1 := t
	t2 := t.Add(time.Minute)
	_, err := db.DbMap.Select(&events, "SELECT * FROM events WHERE CONVERT(?, DATETIME) <= events.startdate and events.startdate < CONVERT(?, DATETIME)", t1, t2)

	if err != nil {
		log.Fatal(err)
	}
	return events
}

func RunQueue() {
	virtualQueue := make([]jobQueue, len(queue))
	copy(virtualQueue, queue)
	queue = make([]jobQueue, 0)

	go RunInstanse(virtualQueue)
}

func RunInstanse(virtualQueue []jobQueue) {
	for _, q := range virtualQueue {
		assistRules, settings, event, configuration, eventRules := confjson.ReadDefaultConfigs()
		confjson.SetConfigs(q.Event, &assistRules, &settings, &event, &configuration, &eventRules)
		// TODO
	}
}
