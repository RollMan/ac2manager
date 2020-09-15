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
	"io"
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

type JobQueue struct {
	JobType JobType
	Event   models.Event
}

type ruleFile struct {
	rule     confjson.Rule
	filename string
}

type OpenCloseWriter interface {
	io.WriteCloser
	OpenFile(string, int, os.FileMode) error
}

type FileOpenCloseWriter struct {
	File *os.File
}

func (f *FileOpenCloseWriter) OpenFile(name string, flag int, perm os.FileMode) error {
	var err error
	f.File, err = os.OpenFile(name, flag, perm)
	return err
}

func (f *FileOpenCloseWriter) Write(p []byte) (int, error) {
	n, err := f.File.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (f *FileOpenCloseWriter) Close() error {
	err := f.File.Close()
	return err
}

type JobmngerAPI interface {
	FindJobs(t time.Time)
	RunQueue()
	RunInstanse(virtualQueue []JobQueue) error
	selectJobsByDate(time.Time) []models.Event
}

type Jobmnger struct {
	Queue       []JobQueue
	DbMap       *gorp.DbMap
	Ec2svc      ec2.Ec2
	DstJsonFile OpenCloseWriter
}

func InitQueue() []JobQueue {
	var queue []JobQueue
	queue = make([]JobQueue, 0)
	return queue
}

func (j *Jobmnger) FindJobs(t time.Time) {
	targetInMinute := t.Truncate(time.Minute)
	events := j.selectJobsByDate(targetInMinute)
	for _, e := range events {
		j.Queue = append(j.Queue, JobQueue{Start, e})
		extra := time.Minute * 10
		enddate := e.Startdate.Add(time.Minute*time.Duration(e.P_sessionDurationMinute+e.Q_sessionDurationMinute+e.R_sessionDurationMinute) + extra)
		end_e := e
		end_e.Startdate = enddate
		j.Queue = append(j.Queue, JobQueue{Stop, end_e})
	}
}

func (j *Jobmnger) selectJobsByDate(t time.Time) []models.Event {
	var events []models.Event
	t1 := t
	t2 := t.Add(time.Minute)
	_, err := j.DbMap.Select(&events, "SELECT * FROM events WHERE CONVERT(?, DATETIME) <= events.startdate and events.startdate < CONVERT(?, DATETIME)", t1, t2)

	if err != nil {
		log.Fatal(err)
	}
	return events
}

func (j *Jobmnger) RunQueue() {
	virtualQueue := make([]JobQueue, len(j.Queue))
	copy(virtualQueue, j.Queue)
	j.Queue = make([]JobQueue, 0)

	go j.RunInstanse(virtualQueue)
}

func (j *Jobmnger) RunInstanse(virtualQueue []JobQueue) error {
	for _, q := range virtualQueue {
		// Select instance to deploy
		// FIXME: create instance from an AMI and to select an available instance.
		id := ""
		if q.JobType == Stop {
			j.Ec2svc.StopInstance(id)
		} else if q.JobType == Start {
			assistRules, settings, event, configuration, eventRules := confjson.ReadDefaultConfigs()
			confjson.SetConfigs(q.Event, assistRules, settings, event, configuration, eventRules)

			conf_dir_path := "/opt/ac2manager/" + id
			err := os.MkdirAll(conf_dir_path, 0777)
			if err != nil {
				j.Ec2svc.StopInstance(id)
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
				conf_path := conf_dir_path + "/" + r.filename
				err := j.DstJsonFile.OpenFile(conf_path, os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					return fmt.Errorf("Failed to open file %s: %v\n", conf_path, err)
				}
				json, err := json.Marshal(r.rule)
				if err != nil {
					j.DstJsonFile.Close()
					return fmt.Errorf("Failed to marshal a json: %s.\n%v", r.filename, r.rule)
				}
				_, err = j.DstJsonFile.Write(json)
				if err != nil {
					j.DstJsonFile.Close()
					return fmt.Errorf("Failed to write a json of %s.\n%v", conf_path, json)
				}
				j.DstJsonFile.Close()
			}

			j.Ec2svc.StartInstance(id)
		}
	}
	return nil
}

var _ JobmngerAPI = &Jobmnger{}
