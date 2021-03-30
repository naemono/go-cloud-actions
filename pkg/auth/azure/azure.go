package azure

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2020-11-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// AuthConfig is the configuration for azure authentication
type AuthConfig struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
	AuxTenantIDs   []string
	Resource       string
}

// NewGroupsClient will return a new azure resource groups client
func NewGroupsClient(conf AuthConfig) (groupsClient resources.GroupsClient, err error) {
	groupsClient = resources.NewGroupsClient(conf.SubscriptionID)

	var a autorest.Authorizer
	a, err = newMgmtAuthorizer(conf)
	if err != nil {
		return groupsClient, errors.Wrap(err, "failed to get new azure groups client")
	}

	groupsClient.Authorizer = a
	groupsClient.AddToUserAgent(fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0"))
	return groupsClient, nil
}

// NewVirtualNetworkPeeringsClient will return a new azure virtual network peerings client
func NewVirtualNetworkPeeringsClient(conf AuthConfig) (vnpc network.VirtualNetworkPeeringsClient, err error) {
	vnpc = network.NewVirtualNetworkPeeringsClient(conf.SubscriptionID)
	if len(conf.AuxTenantIDs) > 0 {
		sender := autorest.CreateSender()
		oauth, err := adal.NewMultiTenantOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, conf.TenantID, conf.AuxTenantIDs, adal.OAuthOptions{})
		if err != nil {
			return vnpc, errors.Wrap(err, "failed to get new ulti-tenant oauth configuration")
		}
		var token *adal.MultiTenantServicePrincipalToken
		token, err = adal.NewMultiTenantServicePrincipalToken(oauth, conf.ClientID, conf.ClientSecret, conf.Resource)
		if err != nil {
			return vnpc, errors.Wrap(err, "failed to generate new multi-tenant service principal token")
		}
		token.PrimaryToken.SetSender(sender)
		for _, t := range token.AuxiliaryTokens {
			t.SetSender(sender)
		}
		vnpc.Authorizer = autorest.NewMultiTenantServicePrincipalTokenAuthorizer(token)
		vnpc.UserAgent = fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0")
		vnpc.Sender = sender
		return vnpc, err
	}
	vnpc.Authorizer, err = auth.NewClientCredentialsConfig(conf.ClientID, conf.ClientSecret, conf.TenantID).Authorizer()
	if err != nil {
		return vnpc, errors.Wrap(err, "failed to authorize with credentials")
	}
	vnpc.UserAgent = fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0")
	vnpc.Sender = autorest.CreateSender()
	return
}

// NewApplicationsClient will return a new azure graph applications client
func NewApplicationsClient(conf AuthConfig) (appsClient graphrbac.ApplicationsClient, err error) {
	appClient := graphrbac.NewApplicationsClient(conf.TenantID)

	var a autorest.Authorizer
	a, err = newAuthorizer(conf)
	if err != nil {
		return appsClient, errors.Wrap(err, "failed to get new azure authorizer")
	}

	appClient.Authorizer = a
	appClient.AddToUserAgent(fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0"))
	return appClient, nil
}

// NewRoleDefinitionsClient will return a new azure role definitions client
func NewRoleDefinitionsClient(conf AuthConfig) (rdClient authorization.RoleDefinitionsClient, err error) {
	rdClient = authorization.NewRoleDefinitionsClient(conf.SubscriptionID)

	var a autorest.Authorizer
	a, err = newMgmtAuthorizer(conf)
	if err != nil {
		return rdClient, errors.Wrap(err, "failed to get new azure role definitions client")
	}
	rdClient.Authorizer = a
	rdClient.AddToUserAgent(fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0"))
	return rdClient, nil
}

// NewServicePrincipalsClient will return a new azure graph service principals client
func NewServicePrincipalsClient(conf AuthConfig) (spClient graphrbac.ServicePrincipalsClient, err error) {
	spClient = graphrbac.NewServicePrincipalsClient(conf.TenantID)
	var a autorest.Authorizer
	a, err = newAuthorizer(conf)
	if err != nil {
		return spClient, errors.Wrap(err, "failed to get new azure authorizer")
	}
	spClient.Authorizer = a
	spClient.AddToUserAgent(fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0"))
	return spClient, nil
}

// NewContainerInstanceClient will return a new azure container groups client
func NewContainerInstanceClient(conf AuthConfig) (cgClient containerinstance.ContainerGroupsClient, err error) {
	cgClient = containerinstance.NewContainerGroupsClient(conf.SubscriptionID)
	var a autorest.Authorizer
	a, err = newMgmtAuthorizer(conf)
	if err != nil {
		return cgClient, errors.Wrap(err, "failed to get new azure authorizer")
	}
	cgClient.Authorizer = a
	cgClient.AddToUserAgent(fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0"))
	return
}

// NewNetworkProfilesClient will return a new azure network profiles client
func NewNetworkProfilesClient(conf AuthConfig) (profClient network.ProfilesClient, err error) {
	profClient = network.NewProfilesClient(conf.SubscriptionID)
	var a autorest.Authorizer
	a, err = newMgmtAuthorizer(conf)
	if err != nil {
		return profClient, errors.Wrap(err, "failed to get new azure authorizer")
	}
	profClient.Authorizer = a
	profClient.AddToUserAgent(fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0"))
	return
}

// NewVirtualNetworksClient will return a new azure virtual networks client
func NewVirtualNetworksClient(conf AuthConfig) (vnetClient network.VirtualNetworksClient, err error) {
	vnetClient = network.NewVirtualNetworksClient(conf.SubscriptionID)
	var a autorest.Authorizer
	a, err = newMgmtAuthorizer(conf)
	if err != nil {
		return vnetClient, errors.Wrap(err, "failed to get new azure authorizer")
	}
	vnetClient.Authorizer = a
	vnetClient.AddToUserAgent(fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0"))
	return
}

// NewSubnetsClient will return a new azure network subnets client
func NewSubnetsClient(conf AuthConfig) (snetClient network.SubnetsClient, err error) {
	snetClient = network.NewSubnetsClient(conf.SubscriptionID)
	var a autorest.Authorizer
	a, err = newMgmtAuthorizer(conf)
	if err != nil {
		return snetClient, errors.Wrap(err, "failed to get new azure authorizer")
	}
	snetClient.Authorizer = a
	snetClient.AddToUserAgent(fmt.Sprintf("Go-Cloud-Actions-v%s", "0.1.0"))
	return
}

func newAuthorizer(conf AuthConfig) (autorest.Authorizer, error) {
	var a autorest.Authorizer

	oauthConfig, err := adal.NewOAuthConfig(
		azure.PublicCloud.ActiveDirectoryEndpoint, conf.TenantID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new azure oauth config")
	}

	token, err := adal.NewServicePrincipalToken(
		*oauthConfig, conf.ClientID, conf.ClientSecret, azure.PublicCloud.GraphEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new azure service principal token")
	}
	a = autorest.NewBearerAuthorizer(token)
	return a, nil
}

func newMgmtAuthorizer(conf AuthConfig) (autorest.Authorizer, error) {
	var a autorest.Authorizer

	oauthConfig, err := adal.NewOAuthConfig(
		azure.PublicCloud.ActiveDirectoryEndpoint, conf.TenantID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new azure oauth config")
	}

	token, err := adal.NewServicePrincipalToken(
		*oauthConfig, conf.ClientID, conf.ClientSecret, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new azure service principal token")
	}
	a = autorest.NewBearerAuthorizer(token)
	return a, nil
}
