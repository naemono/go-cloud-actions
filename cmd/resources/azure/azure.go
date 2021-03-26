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
	azure_resources "github.com/naemono/go-cloud-actions/pkg/resources/azure"
)

var (
	AzureCmd = &cobra.Command{
		Use:              "azure",
		Short:            "Control resources in azure's public clouds",
		Long:             `A cli to interact with resources in Azure's public cloud.`,
		PersistentPreRun: shared_azure.PersistentPreRun,
	}
	resourceGroupsCmd = &cobra.Command{
		Use:   "resource-groups",
		Short: "control resources groups in azure's public clouds",
		Long:  `A cli to control resource groups in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
		},
	}
	resourceGroupAddCmd = &cobra.Command{
		Use:   "add",
		Short: "add resource group in azure's public clouds",
		Long:  `A cli to add resource group in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
			viper.BindPFlag("location", cmd.Flags().Lookup("location"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("name") == "" {
				return fmt.Errorf("name cannot be empty")
			}
			if viper.GetString("location") == "" {
				return fmt.Errorf("location cannot be empty")
			}
			return createResourceGroup()
		},
	}
)

func init() {
	shared_azure.AddAuthFlagsToCommand(AzureCmd)

	resourceGroupAddCmd.Flags().StringP("name", "n", "", "name of resource group")
	resourceGroupAddCmd.Flags().StringP("location", "L", "", "location/region of resource group")

	AzureCmd.AddCommand(resourceGroupsCmd)
	resourceGroupsCmd.AddCommand(resourceGroupAddCmd)
}

func createResourceGroup() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("creating resource group")
	client, err := azure_resources.New(azure_resources.Config{
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
	err = client.CreateResourceGroup(ctx, viper.GetString("name"), strings.ToLower(viper.GetString("location")))
	if err != nil {
		return err
	}
	logger.Infof("resource group '%s' created", viper.GetString("name"))
	return nil
}
