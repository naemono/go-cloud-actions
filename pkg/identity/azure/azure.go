package azure

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/authorization/mgmt/authorization"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"

	azure_auth "github.com/naemono/go-cloud-actions/pkg/auth/azure"
)

var (
	// ErrApplicationAlreadyExists is the error when an application already exists
	ErrApplicationAlreadyExists = errors.New("application already exists")
	// ErrServicePrincipalAlreadyExists is the error when a azure service principal already exists
	ErrServicePrincipalAlreadyExists = errors.New("service principal already exists")
)

// Config is the configuration for the azure users client
type Config struct {
	azure_auth.AuthConfig
}

// Client is the client for the azure users client
type Client struct {
	Config
}

// ApplicationConfig is the configuration for an azure application
type ApplicationConfig struct {
	AppID                   *string
	AvailableToOtherTenants bool
	DisplayName             string
	HomePage                string
	IdentifierUris          []string
}

// New will return a new azure identities client
func New(conf Config) *Client {
	return &Client{
		Config: conf,
	}
}

// CreateADApplication creates an Azure Active Directory (AAD) application
func (c *Client) CreateADApplication(ctx context.Context, appConfig ApplicationConfig) (graphrbac.Application, error) {
	appClient, err := azure_auth.NewApplicationsClient(c.AuthConfig)
	if err != nil {
		return graphrbac.Application{}, errors.Wrap(err, "failed to get new azure applications client")
	}
	var res graphrbac.ApplicationListResultPage
	res, err = appClient.List(ctx, fmt.Sprintf("displayName eq '%s'", appConfig.DisplayName))
	if err != nil {
		return graphrbac.Application{}, errors.Wrap(err, "failed to list applications")
	}
	if len(res.Values()) > 0 {
		return graphrbac.Application{}, ErrApplicationAlreadyExists
	}
	return appClient.Create(ctx, graphrbac.ApplicationCreateParameters{
		AvailableToOtherTenants: &appConfig.AvailableToOtherTenants,
		DisplayName:             &appConfig.DisplayName,
		Homepage:                &appConfig.HomePage,
		IdentifierUris:          &appConfig.IdentifierUris,
	})
}

// CreateServicePrincipal creates a service principal associated with the specified application.
func (c *Client) CreateServicePrincipal(ctx context.Context, appConfig ApplicationConfig) (graphrbac.ServicePrincipal, error) {
	if appConfig.AppID == nil {
		return graphrbac.ServicePrincipal{}, fmt.Errorf("app id cannot be empty")
	}
	spClient, err := azure_auth.NewServicePrincipalsClient(c.AuthConfig)
	if err != nil {
		return graphrbac.ServicePrincipal{}, errors.Wrap(err, "failed to get new azure service principal client")
	}
	var res graphrbac.ServicePrincipalListResultPage
	res, err = spClient.List(ctx, fmt.Sprintf("displayname eq '%s' and servicePrincipalType eq 'Application'", appConfig.DisplayName))
	if err != nil {
		return graphrbac.ServicePrincipal{}, errors.Wrap(err, "failed to list service principals")
	}
	if len(res.Values()) > 0 {
		return graphrbac.ServicePrincipal{}, ErrServicePrincipalAlreadyExists
	}
	return spClient.Create(ctx,
		graphrbac.ServicePrincipalCreateParameters{
			AppID:          appConfig.AppID,
			AccountEnabled: to.BoolPtr(true),
		})
}

// ListRoleDefinitions will list azure role definitions for a given resource group, and virtual network
func (c *Client) ListRoleDefinitions(ctx context.Context, rg, vnet string) ([]authorization.RoleDefinition, error) {
	rdClient, err := azure_auth.NewRoleDefinitionsClient(c.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get new role definitions client")
	}
	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s", c.SubscriptionID, rg, vnet)
	var res authorization.RoleDefinitionListResultPage
	res, err = rdClient.List(ctx, scope, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list role definitions")
	}
	return res.Values(), nil
}

// CreateApplicationCredentials will create/update an existing application's credentials
func (c *Client) CreateApplicationCredentials(ctx context.Context, appConfig ApplicationConfig) (password string, err error) {
	if appConfig.AppID == nil {
		return "", fmt.Errorf("app id cannot be empty")
	}
	password = randomPassword()
	appClient, err := azure_auth.NewApplicationsClient(c.AuthConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to get new azure applications client")
	}
	_, err = appClient.UpdatePasswordCredentials(ctx, *appConfig.AppID, graphrbac.PasswordCredentialsUpdateParameters{
		Value: &[]graphrbac.PasswordCredential{
			{
				StartDate: &date.Time{Time: time.Now()},
				EndDate:   &date.Time{Time: time.Now().Add(24 * 365 * time.Hour)},
				Value:     to.StringPtr(password),
			},
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to update password credentials for app")
	}
	return
}

func randomPassword() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 16
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}
