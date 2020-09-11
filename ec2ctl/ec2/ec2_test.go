package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
	"time"
)

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
