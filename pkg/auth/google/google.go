package google

import (
	"context"

	"github.com/pkg/errors"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

// AuthConfig is the configuration required to generate any google api client
type AuthConfig struct {
	CredentialsFilePath string
}

// NewNetworkClient will return a new google network client with a given configuration
func NewNetworkClient(ctx context.Context, conf AuthConfig) (*compute.NetworksService, error) {
	svc, err := compute.NewService(ctx, option.WithCredentialsFile(conf.CredentialsFilePath))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new google network service")
	}
	return compute.NewNetworksService(svc), nil
}

// NewContainersClient will return a new google containers (gke) client with a given configuration
func NewContainersClient(ctx context.Context, conf AuthConfig) (*container.ProjectsService, error) {
	svc, err := container.NewService(ctx, option.WithCredentialsFile(conf.CredentialsFilePath))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new google containers service")
	}
	return container.NewProjectsService(svc), nil
}
