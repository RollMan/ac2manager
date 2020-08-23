package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
)

var Sess *session.Session
var Ec2Svc *ec2.EC2

func InitAWS() {
	Sess = session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("ap-northeast-1")}, // Tokyo
		SharedConfigState: session.SharedConfigEnable,
	}))

	Ec2Svc = ec2.New(Sess)
}

func StartInstance(instanceId string) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		DryRun: aws.Bool(true),
	}
	result, err := Ec2Svc.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		result, err = Ec2Svc.StartInstances(input)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(result.StartingInstances)
		}
	} else {
		log.Println(err)
	}
}

func DescribeInstance(instanceId string) []*ec2.Reservation {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}
	result, err := Ec2Svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatalln(aerr.Error())
		} else {
			log.Fatalln(err)
		}
	}

	return result.Reservations
}

func DescribeInstanceStatus(instanceId string) []*ec2.InstanceStatus {
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
		IncludeAllInstances: aws.Bool(true),
	}
	result, err := Ec2Svc.DescribeInstanceStatus(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatalln(aerr.Error())
		} else {
			log.Fatalln(err)
		}
	}

	return result.InstanceStatuses
}
