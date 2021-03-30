package google

import (
	"context"
	"fmt"
	"time"

	google_auth "github.com/naemono/go-cloud-actions/pkg/auth/google"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
)

// PeeringCommon are the common fields between create/list peering requests
type PeeringCommon struct {
	ProjectID   string
	NetworkName string
	PeeringName string
}

// CreatePeeringRequest is a request to create a peering between 2 google project's networks
type CreatePeeringRequest struct {
	PeeringCommon
	RemoteNetworkName              string
	RemoteProjectName              string
	ExportCustomRoutes             bool
	ExportSubnetRoutesWithPublicIP bool
	ImportCustomRoutes             bool
	ImportSubnetRoutesWithPublicIP bool
}

// ListPeeringRequest is a request to list peerings for a specified project/network/peering name
type ListPeeringRequest struct {
	PeeringCommon
	Region string
}

// Config is an google peering config
type Config struct {
	google_auth.AuthConfig
	Logger *logrus.Entry
}

// Client is an azure peering client
type Client struct {
	Config
	networksServiceClient *compute.NetworksService
}

// New will return a new google peering client
func New(conf Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	networkClient, err := google_auth.NewNetworkClient(ctx, conf.AuthConfig)
	if err != nil {
		return nil, err
	}
	client := &Client{
		Config:                conf,
		networksServiceClient: networkClient,
	}
	if client.Logger == nil {
		client.Logger = logrus.NewEntry(logrus.New())
		client.Logger.Logger.SetLevel(logrus.InfoLevel)
		client.Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return client, nil
}

// CreatePeering will create a google peering between 2 project's networks
func (c *Client) CreatePeering(ctx context.Context, req CreatePeeringRequest) error {
	remoteNetworkURL := fmt.Sprintf(
		"https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s",
		req.RemoteProjectName,
		req.RemoteNetworkName,
	)
	_, err := c.networksServiceClient.AddPeering(req.ProjectID, req.NetworkName, &compute.NetworksAddPeeringRequest{
		NetworkPeering: &compute.NetworkPeering{
			ExchangeSubnetRoutes:           true,
			ExportCustomRoutes:             req.ExportCustomRoutes,
			ExportSubnetRoutesWithPublicIp: req.ExportSubnetRoutesWithPublicIP,
			ImportCustomRoutes:             req.ImportCustomRoutes,
			ImportSubnetRoutesWithPublicIp: req.ImportSubnetRoutesWithPublicIP,
			Name:                           req.PeeringName,
			Network:                        remoteNetworkURL,
		},
	}).Do()
	if err != nil {
		return errors.Wrap(err, "failed to create peer")
	}
	c.Logger.Infof("peering created succesfully")
	return nil
}

// ListPeerings will list a google project's network peering
func (c *Client) ListPeerings(ctx context.Context, req ListPeeringRequest) error {
	var (
		err error
		res *compute.ExchangedPeeringRoutesList
	)
	for _, direction := range []string{"OUTGOING", "INCOMING"} {
		res, err = c.networksServiceClient.ListPeeringRoutes(req.ProjectID, req.NetworkName).PeeringName(req.PeeringName).Region(req.Region).Direction(direction).Do()
		if err != nil {
			return err
		}
		for _, r := range res.Items {
			c.Logger.Infof("peer: %+v", r)
		}
	}
	return nil
}
