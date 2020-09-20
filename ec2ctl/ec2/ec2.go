package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"log"
)

type Ec2 struct {
	Svc ec2iface.EC2API
}

func InitAWS() Ec2 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("ap-northeast-1")}, // Tokyo
		SharedConfigState: session.SharedConfigEnable,
	}))

	Svc := ec2.New(sess)

	ec2svc := Ec2{Svc}
	return ec2svc
}

// TODO: return error when fail
func (ec2svc *Ec2) StartInstance(instanceId string) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		DryRun: aws.Bool(true),
	}
	result, err := ec2svc.Svc.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = ec2svc.Svc.StartInstances(input)
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Printf("EC2 instance (id: %s) launched\n", instanceId)
			log.Println(result.StartingInstances)
		}
	} else {
		log.Fatalln(err)
	}
}

// TODO: return error when fail
func (ec2svc *Ec2) DescribeInstance(instanceId string) *ec2.Instance {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}
	result, err := ec2svc.Svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatalln(aerr.Error())
		} else {
			log.Fatalln(err)
		}
	}

	var instances []*ec2.Instance

	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			instances = append(instances, i)
		}
	}
	if len(instances) != 1 {
		log.Fatalln("There are more than one instances even though the instance is specified by ID.")
	}
	return instances[0]
}

// TODO: return error when fail
func (ec2svc *Ec2) StopInstance(instanceId string) {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		DryRun: aws.Bool(true),
	}
	result, err := ec2svc.Svc.StopInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = ec2svc.Svc.StopInstances(input)
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Printf("EC2 instance (id: %s) stopped\n", instanceId)
			log.Println(result.StoppingInstances)
		}
	} else {
		log.Fatalln(err)
	}
}

// TODO: return error when fail
func (ec2svc *Ec2) DescribeInstanceStatus(instanceId string) []*ec2.InstanceStatus {
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		IncludeAllInstances: aws.Bool(true),
	}
	result, err := ec2svc.Svc.DescribeInstanceStatus(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatalln(aerr.Error())
		} else {
			log.Fatalln(err)
		}
	}

	if len(result.InstanceStatuses) != 1 {
		log.Fatalln("The number of reservations is not 1 even though the instance is specified by ID.")
	}

	return result.InstanceStatuses
}

func (ec2svc *Ec2) DescribeInstanceIPAddress(instanceId string) (string, error) {
	instance := ec2svc.DescribeInstance(instanceId)
	address_p := instance.PublicIpAddress
	if address_p == nil {
		return "", fmt.Errorf("No IP addresses. Is the instance running?")
	}
	address := *address_p
	return address, nil
}
