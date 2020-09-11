package ec2

import (
	_ "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"testing"
)

type mockedInstance struct {
	ec2iface.EC2API
	Resp ec2.StartInstancesOutput
}

func (m mockedInstance) StartInstances(i *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	// TODO: FIXME: return error when dry run
	if *i.DryRun == true {
		return nil, awserr.New("DryRunOperation", "", nil)
	}
	return &m.Resp, nil
}

func TestStartInstance(t *testing.T) {
	cases := []struct {
		Resp     ec2.StartInstancesOutput
		Expected []ec2.InstanceState
	}{
		{
			Resp: ec2.StartInstancesOutput{
				StartingInstances: []*ec2.InstanceStateChange{
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
		ec2svc := Ec2{
			svc: mockedInstance{Resp: c.Resp},
		}
		ec2svc.StartInstance("dummyinstanceid")
	}
}

/* Danger test due to billing occuring
func TestStartInstance(t *testing.T) {
	const instanceId = "i-0fbc76811e301234a"
	InitAWS()
	StartInstance(instanceId)

	cnt := 0
	max := 20
	for {
		instanceStatus := DescribeInstanceStatus(instanceId)
		if len(instanceStatus) == 0 {
			t.Fatalf("The length of reservation is 0.\n")
		}
		status := instanceStatus[0].InstanceState.Name
		if *status == ec2.InstanceStateNameRunning {
			t.Log("ok")
			break
		} else if *status == ec2.InstanceStateNamePending {
			t.Logf("Still pending. Re-fetching status... (%d/%d)", cnt, max)
			time.Sleep(3 * time.Second)
		} else {
			t.Fatalf("Failed to start instance: %v", status)
		}
	}
}
*/
