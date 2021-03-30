package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerinstance/mgmt/containerinstance"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	azure_auth "github.com/naemono/go-cloud-actions/pkg/auth/azure"
)

// Config is the configuration for the azure serverless package
type Config struct {
	azure_auth.AuthConfig
	Logger *logrus.Entry
}

// Client is the client fo rthe azure serverless package
type Client struct {
	Config
	cgClient containerinstance.ContainerGroupsClient
}

// CreateContainerRequest is a request to create an azure Container Instance
type CreateContainerRequest struct {
	APIVersion               string                                     `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Type                     string                                     `json:"type,omitempty" yaml:"type,omitempty"`
	ContainerGroupName       string                                     `json:"name" yaml:"name"`
	Location                 string                                     `json:"location" yaml:"location"`
	ResourceGroupName        string                                     `json:"resourceGroup" yaml:"resourceGroup"`
	ContainerGroupProperties containerinstance.ContainerGroupProperties `json:"properties" yaml:"properties"`
}

// New will return a new azure serverless (container instances) client
func New(conf Config) (*Client, error) {
	var err error
	c := &Client{
		Config: conf,
	}
	c.cgClient, err = azure_auth.NewContainerInstanceClient(conf.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new container service client")
	}
	if c.Logger == nil {
		c.Logger = logrus.NewEntry(logrus.New())
		c.Logger.Logger.SetLevel(logrus.InfoLevel)
		c.Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return c, nil
}

// CreateContainerGroup creates a new container group given a container group name, location,
// resource group, and container properties in request
func (c *Client) CreateContainerGroup(ctx context.Context, req CreateContainerRequest) (cg containerinstance.ContainerGroup, err error) {
	future, err := c.cgClient.CreateOrUpdate(
		ctx,
		req.ResourceGroupName,
		req.ContainerGroupName,
		containerinstance.ContainerGroup{
			Name:                     &req.ContainerGroupName,
			Location:                 &req.Location,
			ContainerGroupProperties: &req.ContainerGroupProperties,
		})

	if err != nil {
		return cg, errors.Wrap(err, "failed to create container group")
	}

	err = future.WaitForCompletionRef(ctx, c.cgClient.Client)
	if err != nil {
		return cg, errors.Wrap(err, "failed waiting for container group to come online")
	}
	return future.Result(c.cgClient)
}
