package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	azure_auth "github.com/naemono/go-cloud-actions/pkg/auth/azure"
)

// Config is the configuration for the azure resources package
type Config struct {
	azure_auth.AuthConfig
	Logger *logrus.Entry
}

// Client is the client for the azure resources package
type Client struct {
	Config
	groupsClient resources.GroupsClient
}

// New will return a new azure resources client
func New(conf Config) (*Client, error) {
	var err error
	c := &Client{
		Config: conf,
	}
	c.groupsClient, err = azure_auth.NewGroupsClient(conf.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new groups client")
	}
	if c.Logger == nil {
		c.Logger = logrus.NewEntry(logrus.New())
		c.Logger.Logger.SetLevel(logrus.InfoLevel)
		c.Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return c, nil
}

// CreateResourceGroup will create an azure resource group
func (c *Client) CreateResourceGroup(ctx context.Context, name, location string) error {
	_, err := c.groupsClient.CreateOrUpdate(ctx, name, resources.Group{
		Name:     &name,
		Location: &location,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create resource group %s", name)
	}
	return nil
}
