package jobmng

import (
	"encoding/json"
	"fmt"
	"github.com/RollMan/ac2manager/app/models"
	"github.com/RollMan/ac2manager/ec2ctl/confjson"
	_ "github.com/RollMan/ac2manager/ec2ctl/db"
	"github.com/RollMan/ac2manager/ec2ctl/ec2"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
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

type ruleFile struct {
	rule     confjson.Rule
	filename string
}

func InitQueue() []jobQueue {
	var queue []jobQueue
	queue = make([]jobQueue, 0)
	return queue
}

func FindJobs(t time.Time, queue []jobQueue, dbMap *gorp.DbMap) []jobQueue {
	targetInMinute := t.Truncate(time.Minute)
	events := selectJobsByDate(targetInMinute, dbMap)
	for _, e := range events {
		queue = append(queue, jobQueue{Start, e, e.Startdate})
		extra := time.Minute * 10
		enddate := e.Startdate.Add(time.Minute*time.Duration(e.P_sessionDurationMinute+e.Q_sessionDurationMinute+e.R_sessionDurationMinute) + extra)
		queue = append(queue, jobQueue{Stop, e, enddate})
	}
	return queue
}

func selectJobsByDate(t time.Time, dbMap *gorp.DbMap) []models.Event {
	var events []models.Event
	t1 := t
	t2 := t.Add(time.Minute)
	_, err := dbMap.Select(&events, "SELECT * FROM events WHERE CONVERT(?, DATETIME) <= events.startdate and events.startdate < CONVERT(?, DATETIME)", t1, t2)

	if err != nil {
		log.Fatal(err)
	}
	return events
}

func RunQueue(queue []jobQueue, ec2svc ec2.Ec2) []jobQueue {
	virtualQueue := make([]jobQueue, len(queue))
	copy(virtualQueue, queue)
	queue = make([]jobQueue, 0)

	go RunInstanse(virtualQueue, ec2svc)
	return queue
}

func RunInstanse(virtualQueue []jobQueue, ec2svc ec2.Ec2) error {
	for _, q := range virtualQueue {
		// Select instance to deploy
		// FIXME: create instance from an AMI and to select an available instance.
		id := ""
		if q.JobType == Stop {
			ec2svc.StopInstance(id)
		} else if q.JobType == Start {
			assistRules, settings, event, configuration, eventRules := confjson.ReadDefaultConfigs()
			confjson.SetConfigs(q.Event, assistRules, settings, event, configuration, eventRules)

			conf_dir_path := "/opt/ac2manager/" + id
			err := os.MkdirAll(conf_dir_path, 0777)
			if err != nil {
				ec2svc.StopInstance(id)
				return fmt.Errorf("Failed to create directory: %s", conf_dir_path)
			}

			rules := []ruleFile{
				ruleFile{assistRules, "assistRules.json"},
				ruleFile{settings, "settings.json"},
				ruleFile{configuration, "configuration.json"},
				ruleFile{event, "event.json"},
				ruleFile{eventRules, "eventRules.json"},
			}

			for _, r := range rules {
				json, err := json.Marshal(r.rule)
				if err != nil {
					return fmt.Errorf("Failed to marshal a json: %s.\n%v", r.filename, r.rule)
				}
				conf_path := conf_dir_path + "/" + r.filename
				err = ioutil.WriteFile(conf_path, json, 0644)
				if err != nil {
					return fmt.Errorf("Failed to write a json of %s.\n%v", conf_path, json)
				}
			}

			ec2svc.StartInstance(id)
		}
	}
	return nil
}
