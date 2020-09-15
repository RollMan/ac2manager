package main

import (
	"fmt"
	"github.com/RollMan/ac2manager/ec2ctl/db"
	"github.com/RollMan/ac2manager/ec2ctl/ec2"
	"github.com/RollMan/ac2manager/ec2ctl/jobmng"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
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

	var dbMap *gorp.DbMap
	{
		dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/ac2?charset=utf8&parseTime=true", os.Getenv("AC2_DB_USERNAME"), os.Getenv("MYSQL_ROOT_PASSWORD"))
		_, dbMap = db.InitDB(dsn)
	}

	var jobmnger = jobmng.Jobmnger{
		Queue:       jobmng.InitQueue(),
		DbMap:       dbMap,
		Ec2svc:      ec2.InitAWS(),
		DstJsonFile: &jobmng.FileOpenCloseWriter{},
	}
	fmt.Println(jobmnger)
	// TODO: graceful shutdown when SIGINT
	cron(jobmnger)
}

func cron(jobmnger jobmng.Jobmnger) {
	prev := time.Now()
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
			jobmnger.FindJobs(prev)
		} else {
			now_unixminute := int(now.Unix() / 60)

			for {
				prev = prev.Add(time.Minute)
				jobmnger.FindJobs(prev)
				prev_unixminute := int(prev.Unix() / 60)
				if !(prev_unixminute < now_unixminute) {
					break
				}
			}
		}
		jobmnger.RunQueue()
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
