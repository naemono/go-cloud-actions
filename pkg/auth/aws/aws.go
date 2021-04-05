package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/pkg/errors"
)

// AuthConfig is the authentication configuration for aws
type AuthConfig struct {
	Profile string
	Region  string
}

// NewEc2Client will return a new configured ec2 client
func NewEc2Client(auth AuthConfig) (*ec2.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(auth.Region),
		config.WithSharedConfigProfile(auth.Profile),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ec2 credentials provider from credentials")
	}
	return ec2.NewFromConfig(cfg), nil
}
