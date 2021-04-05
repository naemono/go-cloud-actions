package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/logging"
	aws_auth "github.com/naemono/go-cloud-actions/pkg/auth/aws"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Config is an aws network config
type Config struct {
	aws_auth.AuthConfig
	Logger *logrus.Entry
}

// Client is an aws network client
type Client struct {
	Config
	ec2Client *ec2.Client
}

// CreateVpcRequest is a request to create an aws vpc
type CreateVpcRequest struct {
	CidrBlock         string
	InstanceTenacy    types.Tenancy
	TagSpecifications []types.TagSpecification
	DryRun            bool
}

// CreateVpcSubnetRequest is a request to create a subnet within a vpc
type CreateVpcSubnetRequest struct {
	CidrBlock         string
	VPCId             string
	AvailabilityZone  string
	TagSpecifications []types.TagSpecification
	DryRun            bool
}

// ec2Logger is a Logger implementation that wraps the standard library logger, and delegates logging to it's
// Printf method.
type ec2Logger struct {
	Logger *logrus.Entry
}

// Logf logs the given classification and message to the underlying logger.
func (s ec2Logger) Logf(classification logging.Classification, format string, v ...interface{}) {
	if len(classification) != 0 {
		format = string(classification) + " " + format
	}

	s.Logger.Printf(format, v...)
}

// newEc2Logger returns a new ec2Logger
func newEc2Logger(entry *logrus.Entry) ec2Logger {
	return ec2Logger{
		Logger: entry,
	}
}

func withLogger(logger ec2Logger) func(*ec2.Options) {
	return func(o *ec2.Options) {
		o.Logger = logger
	}
}

// New will return a new aws network client
func New(conf Config) (*Client, error) {
	c := &Client{
		Config: conf,
	}
	if c.Logger == nil {
		c.Logger = logrus.NewEntry(logrus.New())
		c.Logger.Logger.SetLevel(logrus.InfoLevel)
		c.Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	var err error
	c.ec2Client, err = aws_auth.NewEc2Client(conf.AuthConfig)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// CreateVPC will create an aws vpc
func (c *Client) CreateVPC(ctx context.Context, request CreateVpcRequest) error {
	if len(request.TagSpecifications) == 0 {
		return fmt.Errorf("tags are required")
	}
	var found bool
	for _, s := range request.TagSpecifications {
		for _, t := range s.Tags {
			if t.Key != nil && strings.ToLower(*t.Key) == "name" {
				found = true
			}
		}
	}
	if !found {
		return fmt.Errorf("name tag is required")
	}
	c.ec2Client.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock:         &request.CidrBlock,
		InstanceTenancy:   request.InstanceTenacy,
		DryRun:            request.DryRun,
		TagSpecifications: request.TagSpecifications,
	}, withLogger(newEc2Logger(c.Logger)))
	return nil
}

// ListVPCs will list vpcs in the region in which the client is configured
func (c *Client) ListVPCs(ctx context.Context) ([]types.Vpc, error) {
	response, err := c.ec2Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		DryRun: false,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list vpcs")
	}
	return response.Vpcs, nil
}

// DeleteVPC will delete the given vpc id
func (c *Client) DeleteVPC(ctx context.Context, id string) error {
	_, err := c.ec2Client.DeleteVpc(ctx, &ec2.DeleteVpcInput{
		VpcId: to.StringPtr(id),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to delete vpc id %s", id)
	}
	c.Logger.Infof("vpc id %s deleted", id)
	return nil
}

// CreateSubnetInVPC will attempt to create a subnet within a given vpc id
func (c *Client) CreateSubnetInVPC(ctx context.Context, request CreateVpcSubnetRequest) error {
	_, err := c.ec2Client.CreateSubnet(ctx, &ec2.CreateSubnetInput{
		CidrBlock:        &request.CidrBlock,
		VpcId:            &request.VPCId,
		AvailabilityZone: &request.AvailabilityZone,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create subnet in vpc")
	}
	return nil
}

// ListSubnetsInVPC will list the existing subnet within a given vpc
func (c *Client) ListSubnetsInVPC(ctx context.Context, vpcID string) ([]types.Subnet, error) {
	out, err := c.ec2Client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   to.StringPtr("vpc-id"),
				Values: []string{vpcID},
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list subnets in vpc")
	}
	return out.Subnets, nil
}
