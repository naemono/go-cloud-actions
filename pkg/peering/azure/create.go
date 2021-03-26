package azure

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
)

// CreatePeeringRequest is a request to creating a peering
type CreatePeeringRequest struct {
	SourceResourceGroup       string
	SourceVnetName            string
	SourcePeeringName         string
	RemoteVnetID              string
	AllowVirtualNetworkAccess bool
	AllowForwardedTraffic     bool
	AllowGatewayTransit       bool
	UseRemoteGateways         bool
}

// Create will create a peering connection
func (c *Client) Create(ctx context.Context, request CreatePeeringRequest) error {
	vnetPeering := network.VirtualNetworkPeering{
		Name:                                  &request.SourcePeeringName,
		VirtualNetworkPeeringPropertiesFormat: getVirtualNetworkPeeringProperties(request),
	}
	c.Logger.Warnf("attempting to create peering with request %+v, and config: %+v", request, c.Config)
	result, err := c.vnpClient.CreateOrUpdate(
		ctx,
		request.SourceResourceGroup,
		request.SourceVnetName,
		request.SourcePeeringName,
		vnetPeering)
	if err != nil {
		return errors.Wrap(err, "unable to create peering")
	}
	resp := result.Response()
	if !(resp.StatusCode < 300) {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("invalid status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("invalid status code: %d, body: %s", resp.StatusCode, string(b))
	}
	c.Logger.Debug("succesfully created peering")
	return nil
}

func getVirtualNetworkPeeringProperties(request CreatePeeringRequest) *network.VirtualNetworkPeeringPropertiesFormat {
	return &network.VirtualNetworkPeeringPropertiesFormat{
		AllowVirtualNetworkAccess: &request.AllowVirtualNetworkAccess,
		AllowForwardedTraffic:     &request.AllowForwardedTraffic,
		AllowGatewayTransit:       &request.AllowGatewayTransit,
		UseRemoteGateways:         &request.UseRemoteGateways,
		RemoteVirtualNetwork: &network.SubResource{
			ID: &request.RemoteVnetID,
		},
	}
}
