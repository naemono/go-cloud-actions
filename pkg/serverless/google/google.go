package google

import (
	"context"
	"fmt"
	"time"

	google_auth "github.com/naemono/go-cloud-actions/pkg/auth/google"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/container/v1"
)

// ClusterCommon are the common fields between cluster operations
type ClusterCommon struct {
	ProjectID   string
	NetworkName string
}

// CreateClusterRequest is a request to create a gke cluster
type CreateClusterRequest struct {
	ClusterCommon
	ClusterIpv4Cidr string
	Description     string
	Location        string
	Name            string
}

// ListPeeringRequest is a request to list peerings for a specified project/network/peering name
type ListPeeringRequest struct {
	ClusterCommon
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
	containersClient *container.ProjectsService
}

// New will return a new google serverless client
func New(conf Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	containersClient, err := google_auth.NewContainerslient(ctx, conf.AuthConfig)
	if err != nil {
		return nil, err
	}
	client := &Client{
		Config:           conf,
		containersClient: containersClient,
	}
	if client.Logger == nil {
		client.Logger = logrus.NewEntry(logrus.New())
		client.Logger.Logger.SetLevel(logrus.InfoLevel)
		client.Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return client, nil
}

// CreateCluster will create a google gke cluster within a given region
func (c *Client) CreateCluster(ctx context.Context, req CreateClusterRequest) error {
	parent := fmt.Sprintf("projects/%s/locations/%s", req.ProjectID, req.Location)
	_, err := c.containersClient.Locations.Clusters.Create(parent, &container.CreateClusterRequest{
		Cluster: &container.Cluster{
			Autopilot: &container.Autopilot{
				Enabled: true,
			},
			ClusterIpv4Cidr:       req.ClusterIpv4Cidr,
			Description:           req.Description,
			InitialClusterVersion: "latest",
			IpAllocationPolicy:    &container.IPAllocationPolicy{},
			Location:              req.Location,
			Name:                  req.Name,
		},
		Parent: parent,
	}).Do()
	if err != nil {
		return errors.Wrap(err, "failed to create cluster")
	}
	c.Logger.Infof("cluster created succesfully")
	return nil
}
