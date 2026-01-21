package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	computev1 "github.com/varunu28/aws-operators/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func createEc2Instance(ec2Instance *computev1.EC2Instance) (createdInstanceInfo *computev1.CreatedInstanceInfo, err error) {
	l := log.Log.WithName("createEc2Instance")

	l.Info("=== STARTING EC2 INSTANCE CREATION ===",
		"ami", ec2Instance.Spec.AmiID,
		"instanceType", ec2Instance.Spec.InstanceType,
		"region", ec2Instance.Spec.Region)

	// ec2 client for creation of EC2 instance
	ec2Client := awsClient(ec2Instance.Spec.Region)

	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(ec2Instance.Spec.AmiID),
		InstanceType: ec2types.InstanceType(ec2Instance.Spec.InstanceType),
		KeyName:      aws.String(ec2Instance.Spec.KeyPair),
		SubnetId:     aws.String(ec2Instance.Spec.Subnet),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	}

	l.Info("=== CALLING AWS RunInstances API ===")
	result, err := ec2Client.RunInstances(context.TODO(), runInput)
	if err != nil {
		l.Error(err, "Failed to create EC2 instance")
		return nil, fmt.Errorf("failed to create EC2 instance: %w", err)
	}

	if len(result.Instances) == 0 {
		l.Error(nil, "No instances returned in RunInstancesOutput")
		return nil, nil
	}

	instance := result.Instances[0]
	l.Info("=== EC2 INSTANCE CREATED SUCCESSFULLY ===", "instanceId", *instance.InstanceId)

	l.Info("=== WAITING FOR INSTANCE TO BE RUNNING ===")

	runWaiter := ec2.NewInstanceRunningWaiter(ec2Client)
	maxWaitTime := 3 * time.Minute

	err = runWaiter.Wait(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{*instance.InstanceId},
	}, maxWaitTime)
	if err != nil {
		l.Error(err, "Failed to wait for instance to be running")
		return nil, fmt.Errorf("failed to wait for instance to be running: %w", err)
	}

	// Get the public IP address & DNS of the running instance
	l.Info("=== CALLING AWS DescribeInstances API to GET INSTANCE DETAILS ===")
	describeInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{*instance.InstanceId},
	}

	describeResult, err := ec2Client.DescribeInstances(context.TODO(), describeInput)
	if err != nil {
		l.Error(err, "failed to describe EC2 instance")
		return nil, fmt.Errorf("failed to describe EC2 instance: %w", err)
	}

	fmt.Println("Descibe result",
		"public ip", *describeResult.Reservations[0].Instances[0].PublicIpAddress,
		"state", describeResult.Reservations[0].Instances[0].State.Name,
	)

	createdInstance := describeResult.Reservations[0].Instances[0]
	createdInstanceInfo = &computev1.CreatedInstanceInfo{
		InstanceID: *createdInstance.InstanceId,
		State:      string(createdInstance.State.Name),
		PublicIP:   derefString(createdInstance.PublicIpAddress),
		PrivateIP:  derefString(createdInstance.PrivateIpAddress),
		PublicDNS:  derefString(createdInstance.PublicDnsName),
		PrivateDNS: derefString(createdInstance.PrivateDnsName),
	}

	l.Info("=== EC2 INSTACE CREATION COMPLETED ===",
		"instanceID", createdInstanceInfo.InstanceID,
		"state", createdInstanceInfo.State,
		"publicIP", createdInstanceInfo.PublicIP,
	)
	return createdInstanceInfo, nil
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return "<nil>"
}
