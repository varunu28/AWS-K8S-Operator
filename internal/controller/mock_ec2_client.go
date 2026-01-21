package controller

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// MockEC2Client is a mock implementation of EC2API for testing
type MockEC2Client struct {
	RunInstancesFunc       func(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
	DescribeInstancesFunc  func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	TerminateInstancesFunc func(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

func (m *MockEC2Client) RunInstances(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error) {
	if m.RunInstancesFunc != nil {
		return m.RunInstancesFunc(ctx, params, optFns...)
	}
	// Default mock response
	return &ec2.RunInstancesOutput{
		Instances: []ec2types.Instance{
			{
				InstanceId:   aws.String("i-mock123456"),
				State:        &ec2types.InstanceState{Name: ec2types.InstanceStateNamePending},
				InstanceType: ec2types.InstanceTypeT2Micro,
			},
		},
	}, nil
}

func (m *MockEC2Client) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	if m.DescribeInstancesFunc != nil {
		return m.DescribeInstancesFunc(ctx, params, optFns...)
	}
	// Default mock response
	return &ec2.DescribeInstancesOutput{
		Reservations: []ec2types.Reservation{
			{
				Instances: []ec2types.Instance{
					{
						InstanceId:       aws.String("i-mock123456"),
						State:            &ec2types.InstanceState{Name: ec2types.InstanceStateNameRunning},
						InstanceType:     ec2types.InstanceTypeT2Micro,
						PublicIpAddress:  aws.String("1.2.3.4"),
						PrivateIpAddress: aws.String("10.0.0.1"),
						PublicDnsName:    aws.String("ec2-1-2-3-4.compute.amazonaws.com"),
						PrivateDnsName:   aws.String("ip-10-0-0-1.ec2.internal"),
					},
				},
			},
		},
	}, nil
}

func (m *MockEC2Client) TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error) {
	if m.TerminateInstancesFunc != nil {
		return m.TerminateInstancesFunc(ctx, params, optFns...)
	}
	// Default mock response
	return &ec2.TerminateInstancesOutput{
		TerminatingInstances: []ec2types.InstanceStateChange{
			{
				InstanceId:    aws.String("i-mock123456"),
				CurrentState:  &ec2types.InstanceState{Name: ec2types.InstanceStateNameShuttingDown},
				PreviousState: &ec2types.InstanceState{Name: ec2types.InstanceStateNameRunning},
			},
		},
	}, nil
}
