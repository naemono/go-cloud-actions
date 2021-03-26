package azure

import (
	"context"

	"github.com/pkg/errors"
)

// List will list existing peering connections
func (c *Client) List(ctx context.Context, resourceGroup, vnet string) error {
	result, err := c.vnpClient.List(ctx, resourceGroup, vnet)
	if err != nil {
		return errors.Wrap(err, "unable to list peerings")
	}
	resp := result.Response()
	if resp.Value != nil {
		if len(*resp.Value) == 0 {
			c.Logger.Debug("no peerings found")
			return nil
		}
		for i, v := range *resp.Value {
			c.Logger.WithField("name", v.Name).Infof("peer %d", i)
		}
		return nil
	}
	c.Logger.Debug("no peerings found")
	return nil
}
