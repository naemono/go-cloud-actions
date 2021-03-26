package azure

import (
	"context"
	"net"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/davecgh/go-spew/spew"

	azure_auth "github.com/naemono/go-cloud-actions/pkg/auth/azure"
)

const aciDelegationServiceName = "Microsoft.ContainerInstance/containerGroups"

type Config struct {
	azure_auth.AuthConfig
	Logger *logrus.Entry
}

type Client struct {
	Config
	profClient network.ProfilesClient
	vnetClient network.VirtualNetworksClient
	snetClient network.SubnetsClient
}

type NetworkProfileRequest struct {
	Name              string
	ResourceGroupName string
	Location          string
	VnetName          string
	VnetAddressCIDR   string
	SubnetName        string
	SubnetAddressCIDR string
}

func New(conf Config) (*Client, error) {
	var err error
	c := &Client{
		Config: conf,
	}
	c.profClient, err = azure_auth.NewNetworkProfilesClient(conf.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new network profiles client")
	}
	c.vnetClient, err = azure_auth.NewVirtualNetworksClient(conf.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new  virtual networks client")
	}
	c.snetClient, err = azure_auth.NewSubnetsClient(conf.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new network subnets client")
	}
	if c.Logger == nil {
		c.Logger = logrus.NewEntry(logrus.New())
		c.Logger.Logger.SetLevel(logrus.InfoLevel)
		c.Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return c, nil
}

func (c *Client) CreateNetworkProfile(ctx context.Context, req NetworkProfileRequest) error {
	err := validateNetworkProfileRequest(req)
	if err != nil {
		return err
	}
	if req.VnetAddressCIDR == "" {
		req.VnetAddressCIDR = "10.0.0.0/16"
	}
	if req.SubnetAddressCIDR == "" {
		req.SubnetAddressCIDR = "10.0.0.0/24"
	}
	err = c.ensureVnet(ctx, req)
	if err != nil {
		return err
	}
	var snet network.Subnet
	snet, err = c.ensureSubnet(ctx, req)
	if err != nil {
		return err
	}
	// Path properties.containerNetworkInterfaceConfigurations[0].properties.ipConfigurations[0].properties.subnet.
	_, err = c.profClient.CreateOrUpdate(ctx, req.ResourceGroupName, req.Name, network.Profile{
		Location: to.StringPtr(req.Location),
		ProfilePropertiesFormat: &network.ProfilePropertiesFormat{
			ContainerNetworkInterfaceConfigurations: &[]network.ContainerNetworkInterfaceConfiguration{
				{
					Name: to.StringPtr("eth0"),
					ContainerNetworkInterfaceConfigurationPropertiesFormat: &network.ContainerNetworkInterfaceConfigurationPropertiesFormat{
						IPConfigurations: &[]network.IPConfigurationProfile{
							{
								Name: to.StringPtr("ipconfigprofile"),
								IPConfigurationProfilePropertiesFormat: &network.IPConfigurationProfilePropertiesFormat{
									Subnet: &network.Subnet{
										Name: &req.SubnetName,
										ID:   snet.ID,
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create network profile %s", req.Name)
	}
	return nil
}

func (c *Client) ensureVnet(ctx context.Context, req NetworkProfileRequest) (err error) {
	_, err = c.vnetClient.Get(ctx, req.ResourceGroupName, req.VnetName, "")
	if err != nil && strings.Contains(err.Error(), "found") {
		c.Logger.Infof("vnet %s was not found, attempting create", req.VnetName)
		res, err := c.vnetClient.CreateOrUpdate(ctx, req.ResourceGroupName, req.VnetName, network.VirtualNetwork{
			Location: to.StringPtr(req.Location),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: to.StringSlicePtr([]string{req.VnetAddressCIDR}),
				},
			},
		})
		if err != nil {
			return errors.Wrap(err, "failed to create vnet")
		}
		if err = res.WaitForCompletionRef(ctx, c.vnetClient.Client); err != nil {
			return errors.Wrap(err, "failed to wait on vnet creation")
		}
		c.Logger.Infof("vnet %s was created", req.VnetName)
	} else if err != nil {
		return errors.Wrap(err, "request to create vnet failed")
	}
	c.Logger.Infof("vnet %s already exists", req.VnetName)
	return nil
}

func (c *Client) ensureSubnet(ctx context.Context, req NetworkProfileRequest) (snet network.Subnet, err error) {
	snet, err = c.snetClient.Get(ctx, req.ResourceGroupName, req.VnetName, req.SubnetName, "")
	if err != nil && strings.Contains(err.Error(), "found") {
		c.Logger.Infof("subnet %s was not found, attempting create", req.SubnetName)
		res, err := c.snetClient.CreateOrUpdate(ctx, req.ResourceGroupName, req.VnetName, req.SubnetName, network.Subnet{
			Name: &req.SubnetName,
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix: &req.SubnetAddressCIDR,
				Delegations: &[]network.Delegation{
					{
						Name: to.StringPtr(aciDelegationServiceName),
						ServiceDelegationPropertiesFormat: &network.ServiceDelegationPropertiesFormat{
							ServiceName: to.StringPtr(aciDelegationServiceName),
						},
					},
				},
			},
		})
		if err != nil {
			return snet, errors.Wrap(err, "failed to create subnet")
		}
		if err = res.WaitForCompletionRef(ctx, c.snetClient.Client); err != nil {
			return snet, errors.Wrap(err, "failed to wait on subnet creation")
		}
		c.Logger.Infof("subnet %s was created", req.SubnetName)
		return res.Result(c.snetClient)
	} else if err != nil {
		return snet, errors.Wrap(err, "request to create subnet failed")
	}
	c.Logger.Infof("subnet %s already exists", req.SubnetName)
	return snet, nil
}

func validateNetworkProfileRequest(req NetworkProfileRequest) (err error) {
	if req.Location == "" {
		return errors.New("location cannot be empty")
	}
	if req.ResourceGroupName == "" {
		return errors.New("resource group cannot be empty")
	}
	if req.Name == "" {
		return errors.New("name cannot be empty")
	}
	if req.VnetName == "" {
		return errors.New("vnet name cannot be empty")
	}
	if req.VnetAddressCIDR != "" {
		_, _, err = net.ParseCIDR(req.VnetAddressCIDR)
		if err != nil {
			return errors.Wrap(err, "vnet address cidr is invalid")
		}
	}
	if req.SubnetName == "" {
		return errors.New("subnet name cannot be empty")
	}
	if req.SubnetAddressCIDR != "" {
		_, _, err = net.ParseCIDR(req.SubnetAddressCIDR)
		if err != nil {
			return errors.Wrap(err, "subnet address cidr is invalid")
		}
	}
	return nil
}

func (c *Client) ListNetworkProfiles(ctx context.Context, resourceGroupName string) error {
	res, err := c.profClient.List(ctx, resourceGroupName)
	if err != nil {
		return errors.Wrapf(err, "failed to list network profile in resource group %s", resourceGroupName)
	}
	for _, r := range res.Values() {
		c.Logger.Infof("profile: %s", spew.Sdump(r))
	}
	return nil
}
