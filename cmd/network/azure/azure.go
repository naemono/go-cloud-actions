package azure

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	shared_azure "github.com/naemono/go-cloud-actions/cmd/shared/azure"
	auth_azure "github.com/naemono/go-cloud-actions/pkg/auth/azure"
	"github.com/naemono/go-cloud-actions/pkg/logging"
	azure_network "github.com/naemono/go-cloud-actions/pkg/network/azure"
)

var (
	// AzureCmd is the base azure network command
	AzureCmd = &cobra.Command{
		Use:              "azure",
		Short:            "Control networks in azure's public clouds",
		Long:             `A cli to interact with networks in Azure's public cloud.`,
		PersistentPreRun: shared_azure.PersistentPreRun,
	}
	networkProfileCmd = &cobra.Command{
		Use:   "network-profile",
		Short: "control resources groups in azure's public clouds",
		Long:  `A cli to control resource groups in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
		},
	}
	networkProfileAddCmd = &cobra.Command{
		Use:   "add",
		Short: "add network profile in azure's public clouds",
		Long:  `A cli to add network profile in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
			viper.BindPFlag("resource-group", cmd.Flags().Lookup("resource-group"))
			viper.BindPFlag("location", cmd.Flags().Lookup("location"))
			viper.BindPFlag("vnet-name", cmd.Flags().Lookup("vnet-name"))
			viper.BindPFlag("subnet-name", cmd.Flags().Lookup("subnet-name"))
			viper.BindPFlag("vnet-cidr", cmd.Flags().Lookup("vnet-cidr"))
			viper.BindPFlag("subnet-cidr", cmd.Flags().Lookup("subnet-cidr"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("name") == "" {
				return fmt.Errorf("name cannot be empty")
			}
			if viper.GetString("resource-group") == "" {
				return fmt.Errorf("resource-group cannot be empty")
			}
			if viper.GetString("location") == "" {
				return fmt.Errorf("location cannot be empty")
			}
			if viper.GetString("vnet-name") == "" {
				return fmt.Errorf("vnet-name cannot be empty")
			}
			if viper.GetString("subnet-name") == "" {
				return fmt.Errorf("subnet-name cannot be empty")
			}
			return createNetworkProfile()
		},
	}
	networkProfileListCmd = &cobra.Command{
		Use:   "list",
		Short: "list network profile in azure's public clouds",
		Long:  `A cli to list network profile in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("resource-group", cmd.Flags().Lookup("resource-group"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("resource-group") == "" {
				return fmt.Errorf("resource-group cannot be empty")
			}
			return listNetworkProfiles()
		},
	}
)

func init() {
	shared_azure.AddAuthFlagsToCommand(AzureCmd)

	networkProfileAddCmd.Flags().StringP("name", "n", "", "name of network profile")
	networkProfileAddCmd.Flags().StringP("resource-group", "r", "", "name of resource group")
	networkProfileAddCmd.Flags().StringP("location", "L", "", "location/region of network profile")
	networkProfileAddCmd.Flags().StringP("vnet-name", "v", "", "name of the virtual network to use/create")
	networkProfileAddCmd.Flags().StringP("subnet-name", "N", "", "name of the subnet to use/create")
	networkProfileAddCmd.Flags().StringP("vnet-cidr", "V", "10.0.0.0/16", "virtual network cidr to use")
	networkProfileAddCmd.Flags().StringP("subnet-cidr", "C", "10.0.0.0/24", "subnet cidr to use")

	networkProfileListCmd.Flags().StringP("resource-group", "r", "", "name of resource group")

	AzureCmd.AddCommand(networkProfileCmd)
	networkProfileCmd.AddCommand(networkProfileAddCmd)
	networkProfileCmd.AddCommand(networkProfileListCmd)
}

func createNetworkProfile() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("creating network profile")
	client, err := azure_network.New(azure_network.Config{
		AuthConfig: auth_azure.AuthConfig{
			SubscriptionID: viper.GetString("subscription-id"),
			ClientID:       viper.GetString("client-id"),
			ClientSecret:   viper.GetString("client-secret"),
			TenantID:       viper.GetString("tenant-id"),
		},
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	request := azure_network.NetworkProfileRequest{
		Name:              viper.GetString("name"),
		ResourceGroupName: viper.GetString("resource-group"),
		Location:          strings.ToLower(viper.GetString("location")),
		VnetName:          viper.GetString("vnet-name"),
		SubnetName:        viper.GetString("subnet-name"),
	}
	err = client.CreateNetworkProfile(ctx, request)
	if err != nil {
		return err
	}
	logger.Infof("network profile '%s' created", viper.GetString("name"))
	return nil
}

func listNetworkProfiles() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("listing network profiles")
	client, err := azure_network.New(azure_network.Config{
		AuthConfig: auth_azure.AuthConfig{
			SubscriptionID: viper.GetString("subscription-id"),
			ClientID:       viper.GetString("client-id"),
			ClientSecret:   viper.GetString("client-secret"),
			TenantID:       viper.GetString("tenant-id"),
		},
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.ListNetworkProfiles(ctx, viper.GetString("resource-group"))
	if err != nil {
		return err
	}
	return nil
}
