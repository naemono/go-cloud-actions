package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
)

// ListAvailabilityZones will list availability zones available to use
func (c *Client) ListAvailabilityZones(ctx context.Context) ([]types.AvailabilityZone, error) {
	output, err := c.ec2Client.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list availability zones")
	}
	return output.AvailabilityZones, nil
}
