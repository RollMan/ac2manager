package jobmng

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RollMan/ac2manager/app/models"
	ec2svc "github.com/RollMan/ac2manager/ec2ctl/ec2"
	_ "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/go-gorp/gorp"
	"os"
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

	jobmnger := Jobmnger{
		Queue:  nil,
		DbMap:  dbMap,
		Ec2svc: ec2svc.Ec2{},
	}

	target_time := time.Date(2020, 9, 12, 10, 30, 0, 0, time.UTC)
	target_Time2 := target_time.Add(time.Minute)
	{
		row := sqlmock.NewRows([]string{"id", "startdate"}).AddRow(0, target_time)
		expected := models.Event{Id: 0, Startdate: target_time}
		mock.ExpectQuery(`SELECT \* FROM events`).
			WithArgs(target_time, target_Time2).
			WillReturnRows(row)

		events := jobmnger.SelectJobsByDate(target_time)
		if events[0] != expected {
			t.Errorf("invalid result")
		}
	}

	{
		mock.ExpectQuery(`SELECT \* FROM events`).
			WithArgs(target_time, target_Time2).
			WillReturnRows(sqlmock.NewRows([]string{"id", "startdate"}))

		events := jobmnger.SelectJobsByDate(target_time)
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
		jobmnger := Jobmnger{
			Queue:  InitQueue(),
			DbMap:  dbMap,
			Ec2svc: ec2svc.Ec2{},
		}

		mock.ExpectQuery(`SELECT \* FROM events`).
			WithArgs(c.time, c.time.Add(time.Minute)).
			WillReturnRows(c.rows)

		jobmnger.FindJobs(c.time)

		res := jobmnger.Queue

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

type mockedInstanceForTestRunInstance struct {
	ec2iface.EC2API
	RespStart ec2.StartInstancesOutput
	RespStop  ec2.StopInstancesOutput
}

func (m *mockedInstanceForTestRunInstance) StartInstances(i *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	if *i.DryRun == true {
		return nil, awserr.New("DryRunOperation", "", nil)
	}
	return &m.RespStart, nil
}

func (m *mockedInstanceForTestRunInstance) StopInstances(i *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	if *i.DryRun == true {
		return nil, awserr.New("DryRunOperation", "", nil)
	}
	return &m.RespStop, nil
}

type mockedDstJson struct {
	currentName string
	file        map[string][]byte
	opened      bool
}

func (m *mockedDstJson) OpenFile(name string, flag int, perm os.FileMode) error {
	m.currentName = name
	m.opened = true
	return nil
}

func (m *mockedDstJson) Write(p []byte) (int, error) {
	m.file[m.currentName] = p
	return len(p), nil
}

func (m *mockedDstJson) Close() error {
	m.opened = false
	m.currentName = ""
	return nil
}

func TestRunInstance(t *testing.T) {

	type Case struct {
		description  string
		virtualQueue []JobQueue
		RespStart    ec2.StartInstancesOutput
		RespStop     ec2.StopInstancesOutput
	}
	time1 := time.Date(2020, 5, 3, 23, 0, 0, 0, time.UTC)
	// time2 := time.Date(2020, 5, 5, 11, 30, 0, 0, time.UTC)

	cases := []Case{
		{
			description:  "No queue 1",
			virtualQueue: []JobQueue{},
			RespStart:    ec2.StartInstancesOutput{},
			RespStop:     ec2.StopInstancesOutput{},
		},
		{
			description: "Start",
			virtualQueue: []JobQueue{
				{
					JobType: Start,
					Event: models.Event{
						Id:                        0,
						Startdate:                 time1,
						Track:                     "monza_2019",
						WeatherRandomness:         3,
						IsRefuellingAllowedInRace: true,
					},
				},
			},
			RespStart: ec2.StartInstancesOutput{
				StartingInstances: []*ec2.InstanceStateChange{
					{
						CurrentState: &ec2.InstanceState{
							Code: &(&struct{ x int64 }{16}).x,
							Name: &(&struct{ s string }{ec2.InstanceStateNameRunning}).s,
						},
					},
				},
			},
			RespStop: ec2.StopInstancesOutput{},
		},
		{
			description: "Stop",
			virtualQueue: []JobQueue{
				{
					JobType: Stop,
					Event: models.Event{
						Id:           45,
						Startdate:    time1,
						TyreSetCount: 3,
					},
				},
			},
			RespStart: ec2.StartInstancesOutput{},
			RespStop: ec2.StopInstancesOutput{
				StoppingInstances: []*ec2.InstanceStateChange{
					{
						CurrentState: &ec2.InstanceState{
							Code: &(&struct{ x int64 }{16}).x,
							Name: &(&struct{ s string }{ec2.InstanceStateNameRunning}).s,
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		jobmnger := Jobmnger{
			Queue: InitQueue(),
			DbMap: nil,
			Ec2svc: ec2svc.Ec2{
				Svc: &mockedInstanceForTestRunInstance{
					RespStart: c.RespStart,
					RespStop:  c.RespStop,
				},
			},
			DstJsonFile: &mockedDstJson{file: make(map[string][]byte, 0)},
		}

		err := jobmnger.RunInstanse(c.virtualQueue)

		if err != nil {
			t.Errorf("Error in test %s: %v\n", c.description, err)
		}

		value, ok := jobmnger.DstJsonFile.(*mockedDstJson)
		if !ok {
			t.Errorf("Failed type assertion\n")
		}

		for filename, content := range value.file {
			t.Logf("name: %s\ncontent: %s\n", filename, content)
		}
	}
}
