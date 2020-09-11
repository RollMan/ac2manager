package main

import (
	"fmt"
	"github.com/RollMan/ac2manager/ec2ctl/db"
	"github.com/RollMan/ac2manager/ec2ctl/jobmng"
	"log"
	"os"
	"time"
)

type WakeupKind int

const (
	Small = iota
	Medium
	Large
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	time.Local = time.FixedZone("GMT", 0)
	{
		dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/ac2?charset=utf8&parseTime=true", os.Getenv("AC2_DB_USERNAME"), os.Getenv("MYSQL_ROOT_PASSWORD"))
		db.InitDB(dsn)
	}
	jobmng.InitQueue()
	prev := time.Now()
	// TODO: graceful shutdown when SIGINT
	for {
		var now time.Time
		for {
			sleep_by := prev.Add(time.Minute * time.Duration(1)).Truncate(time.Minute)
			sleepUntilNextMinute(sleep_by)

			now = time.Now()
			now_unixminute := int(now.Unix() / 60)
			prev_unixminute := int(now.Unix() / 60)
			if now_unixminute != prev_unixminute {
				break
			}
		}

		timeDiff := now.Sub(prev)
		timeDiffMinute := int(timeDiff.Minutes())

		if timeDiffMinute == 1 {
			prev = now
			jobmng.FindJobs(prev)
		} else {
			now_unixminute := int(now.Unix() / 60)

			for {
				prev = prev.Add(time.Minute)
				jobmng.FindJobs(prev)
				prev_unixminute := int(prev.Unix() / 60)
				if !(prev_unixminute < now_unixminute) {
					break
				}
			}
		}
		jobmng.RunQueue()
	}
}

func sleepUntilNextMinute(target time.Time) {
	t1 := time.Now()
	toWait := target.Sub(t1) + time.Second
	for toWait > 0 {
		time.Sleep(toWait)
		t2 := time.Now()
		toWait = toWait - (t2.Sub(t1))
		t1 = t2
	}
}
