package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	shared_azure "github.com/naemono/go-cloud-actions/cmd/shared/azure"
	auth_azure "github.com/naemono/go-cloud-actions/pkg/auth/azure"
	"github.com/naemono/go-cloud-actions/pkg/logging"
	peering_azure "github.com/naemono/go-cloud-actions/pkg/peering/azure"
	"github.com/naemono/go-cloud-actions/pkg/validate"
)

var (
	// AzureCmd is the base azure peering command
	AzureCmd = &cobra.Command{
		Use:              "azure",
		Short:            "Control peering of VNets in azure's public clouds",
		Long:             `A cli to interact with peering of VNets in Azure's public cloud.`,
		PersistentPreRun: shared_azure.PersistentPreRun,
	}
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create peers of VNets in azure's public clouds",
		Long:  `A cli to create peers of VNets in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("source-peering-name", cmd.Flags().Lookup("source-peering-name"))
			viper.BindPFlag("source-resource-group", cmd.Flags().Lookup("source-resource-group"))
			viper.BindPFlag("source-virtual-network", cmd.Flags().Lookup("source-virtual-network"))
			viper.BindPFlag("target-tenant-id", cmd.Flags().Lookup("target-tenant-id"))
			viper.BindPFlag("target-resource-group", cmd.Flags().Lookup("target-resource-group"))
			viper.BindPFlag("target-virtual-network", cmd.Flags().Lookup("target-virtual-network"))
			viper.BindPFlag("target-subscription-id", cmd.Flags().Lookup("target-subscription-id"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(
				viper.GetViper(),
				[]string{
					"source-resource-group", "source-virtual-network", "source-peering-name", "target-tenant-id",
					"target-resource-group", "target-virtual-network", "target-subscription-id"}); err != nil {
				return err
			}
			return createPeering()
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list peers of VNets in azure's public clouds",
		Long:  `A cli to list peers of VNets in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("resource-group", cmd.Flags().Lookup("resource-group"))
			viper.BindPFlag("vnet-name", cmd.Flags().Lookup("vnet-name"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPeerings()
		},
	}
)

func init() {
	shared_azure.AddAuthFlagsToCommand(AzureCmd)

	createCmd.Flags().StringP("source-resource-group", "r", "", "source resource group where vnet lives")
	createCmd.Flags().StringP("source-virtual-network", "v", "", "source virtual network name within resource group")
	createCmd.Flags().StringP("source-peering-name", "p", "", "source peering name")
	createCmd.Flags().StringP("target-tenant-id", "i", "", "target tenant id (the target tenant where the peering is connecting to)")
	createCmd.Flags().StringP("target-resource-group", "R", "", "target resource group where remote vnet lives")
	createCmd.Flags().StringP("target-virtual-network", "V", "", "target virtual network name within target resource group")
	createCmd.Flags().StringP("target-subscription-id", "T", "", "target subscription id where the target virtual network within target resource group exists")

	listCmd.Flags().StringP("resource-group", "r", "", "resource group in which to list peers")
	listCmd.Flags().StringP("vnet-name", "v", "", "virtual network in which to list peers")

	AzureCmd.AddCommand(createCmd)
	AzureCmd.AddCommand(listCmd)
}

func createPeering() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("creating peering")
	targetResource := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s",
		viper.GetString("target-subscription-id"),
		viper.GetString("target-resource-group"),
		viper.GetString("target-virtual-network"),
	)
	c, err := peering_azure.New(peering_azure.Config{
		AuthConfig: auth_azure.AuthConfig{
			SubscriptionID: viper.GetString("subscription-id"),
			ClientID:       viper.GetString("client-id"),
			ClientSecret:   viper.GetString("client-secret"),
			TenantID:       viper.GetString("tenant-id"),
			AuxTenantIDs: []string{
				viper.GetString("target-tenant-id"),
			},
			Resource: azure.PublicCloud.ResourceManagerEndpoint,
		},
		Logger: logger,
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return c.Create(
		ctx,
		peering_azure.CreatePeeringRequest{
			SourceResourceGroup:       viper.GetString("source-resource-group"),
			SourceVnetName:            viper.GetString("source-virtual-network"),
			SourcePeeringName:         viper.GetString("source-peering-name"),
			RemoteVnetID:              targetResource,
			AllowVirtualNetworkAccess: true,
		},
	)
}

func listPeerings() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	c, err := peering_azure.New(peering_azure.Config{
		AuthConfig: auth_azure.AuthConfig{
			SubscriptionID: viper.GetString("subscription-id"),
			ClientID:       viper.GetString("client-id"),
			ClientSecret:   viper.GetString("client-secret"),
			TenantID:       viper.GetString("tenant-id"),
		},
		Logger: logger,
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return c.List(
		ctx,
		viper.GetString("resource-group"),
		viper.GetString("vnet-name"))
}
