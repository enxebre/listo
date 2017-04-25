package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)

func createInstance(userdata string) (runResult *ec2.Reservation, err error) {
	svc := ec2.New(session.New(&aws.Config{Region: aws.String("eu-west-2")}))
	// Specify the details of the instance that you want to create.
	runResult, err = svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String("ami-16150172"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		SecurityGroups: aws.StringSlice([]string{"listo"}),
		KeyName: aws.String("listo"),
		UserData: &userdata,
	})

	if err != nil {
		log.Println("AWS Client could not create instance", err)
		return
	}

	log.Println("AWS Client created instance", *runResult.Instances[0].InstanceId)

	// Add tags to the created instance
	_ , errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("Listo-Ghost"),
			},
			{
				Key:   aws.String("group"),
				Value: aws.String("listo"),
			},
		},
	})
	if errtag != nil {
		log.Println("AWS Client could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return
	}
	log.Println("AWS Client successfully tagged instance")
	return
}

func getInstances(instanceId string) (ouput *ec2.DescribeInstancesOutput, err error){
	// Only grab instances that are running or just started
	svc := ec2.New(session.New(&aws.Config{Region: aws.String("eu-west-2")}))
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
			{
				Name: aws.String("tag:group"),
				Values: []*string{
					aws.String("listo"),
				},
			},
		},
	}
	if instanceId != "" {
		var  instanceIdFilter = &ec2.Filter {
			Name: aws.String("instance-id"),
			Values: []*string{
				aws.String(instanceId),
			},
		}
		params.Filters = append(params.Filters, instanceIdFilter)
	}
	ouput , err = svc.DescribeInstances(params)
	if err != nil {
		log.Println("AWS Client could not describe instances", err)
		return
	}
	log.Println("AWS Client successfully described instances")
	return
}

func deleteInstance(instanceId string) (ouput *ec2.TerminateInstancesOutput, err error){
	svc := ec2.New(session.New(&aws.Config{Region: aws.String("eu-west-2")}))
	params := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{ // Required
			aws.String(instanceId), // Required
			// More values...
		},
	}
	ouput, err = svc.TerminateInstances(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println("AWS Client could not terminate instance", err)
		return
	}

	log.Println("AWS Client successfully terminated instance")
	return
}
