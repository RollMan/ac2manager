package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RollMan/ac2manager/app/models"
	"github.com/RollMan/ac2manager/ec2ctl/jobmng"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/go-gorp/gorp"
	"os"
	"testing"
	"time"
)

func logTimeOfCalls(title string, times []time.Time, tst *testing.T) {
	tst.Log(title)
	for i, t := range times {
		if i != 0 {
			tst.Logf(" ")
		}
		tst.Logf(t.Format("2006-01-02T15:04:05.000"))
	}
}

type mockedEc2Svc struct {
	ec2iface.EC2API
}

func (m *mockedEc2Svc) StartInstances(i *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	if *i.DryRun == true {
		return nil, awserr.New("DryRunOperation", "", nil)
	}
	return &ec2.StartInstancesOutput{}, nil
}

func (m *mockedEc2Svc) StopInstances(i *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	if *i.DryRun == true {
		return nil, awserr.New("DryRunOperation", "", nil)
	}
	return &ec2.StopInstancesOutput{}, nil
}

type mockedJobmnger00 struct {
	jobmng.Jobmnger
	timeOfFindJobsCall []time.Time
	timeOfRunQueueCall []time.Time
}

func (m *mockedJobmnger00) FindJobs(t time.Time) {
	m.timeOfFindJobsCall = append(m.timeOfFindJobsCall, time.Now())
}
func (m *mockedJobmnger00) RunQueue() {
	m.timeOfRunQueueCall = append(m.timeOfRunQueueCall, time.Now())
}
func (m *mockedJobmnger00) RunInstanse(virtualQueue []jobmng.JobQueue) error {
	return nil
}
func (m *mockedJobmnger00) SelectJobsByDate(time.Time) []models.Event {
	return []models.Event{}
}

func TestCron00(t *testing.T) {
	jobmnger := &mockedJobmnger00{}
	prev := time.Now()
	for i := 0; i < 3; i++ {
		t.Log(i)
		prev = cron(jobmnger, prev)
	}

	logTimeOfCalls("FindJobs", jobmnger.timeOfFindJobsCall, t)
	logTimeOfCalls("FindJobs", jobmnger.timeOfRunQueueCall, t)
}

type mockedJobmnger01 struct {
	jobmng.Jobmnger
	timeOfRunInstanceCall []time.Time
}

type mockedDstJson struct {
	WriteTimes []time.Time
}

func (m *mockedDstJson) OpenFile(name string, flag int, perm os.FileMode) error {
	return nil
}

func (m *mockedDstJson) Write(p []byte) (int, error) {
	m.WriteTimes = append(m.WriteTimes, time.Now())
	return 0, nil
}

func (m *mockedDstJson) Close() error {
	return nil
}

func TestCron01(t *testing.T) {
	const fmt = "2006-01-02T15:04:05"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error %s occuered when opening db mock", err)
	}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbMap.AddTableWithName(models.Event{}, "events").SetKeys(true, "id")
	jobmnger := &mockedJobmnger01{}
	jobmnger.DbMap = dbMap
	jobmnger.DstJsonFile = &mockedDstJson{}
	jobmnger.Ec2svc.Svc = &mockedEc2Svc{}
	defer dbMap.Db.Close()

	now := time.Now()
	target_time := now.Add(5 * time.Minute).Truncate(time.Minute)
	target_time2 := target_time.Add(time.Minute)
	end := target_time.Add(time.Minute + 5*time.Second)

	emptyRow := sqlmock.NewRows([]string{"id", "startdate"})
	row := sqlmock.NewRows([]string{"id", "startdate"}).AddRow(123, target_time)
	for i := 0; i < 4; i++ {
		mock.ExpectQuery(`SELECT \* FROM events`).
			WillReturnRows(emptyRow)
	}
	mock.ExpectQuery(`SELECT \* FROM events`).
		WithArgs(target_time, target_time2).
		WillReturnRows(row)

	mock.ExpectQuery(`SELECT \* FROM events`).
		WillReturnRows(emptyRow)
	mock.ExpectQuery(`SELECT \* FROM events`).
		WillReturnRows(emptyRow)

	cnt := 0
	prev := time.Now()
	for {
		prev = cron(jobmnger, prev)
		t.Logf("loop: %d, prev: %s, end: %s\n", cnt, prev.Format(fmt), end.Format(fmt))
		cnt++
		if time.Now().After(end) {
			break
		}
	}

	wt := jobmnger.DstJsonFile.(*mockedDstJson).WriteTimes[0]
	if (target_time.Add(-15 * time.Second)).Before(wt) && wt.Before(target_time.Add(15*time.Second)) {
	} else {
		t.Errorf("The difference of expected startinstance time is too big:\ntarget:%s, result:%s",
			target_time.Format(fmt), wt.Format(fmt))
	}

	t.Log(wt.Format(fmt))

}
