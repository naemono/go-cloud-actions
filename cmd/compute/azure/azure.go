package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerinstance/mgmt/containerinstance"

	shared_azure "github.com/naemono/go-cloud-actions/cmd/shared/azure"
	auth_azure "github.com/naemono/go-cloud-actions/pkg/auth/azure"
	"github.com/naemono/go-cloud-actions/pkg/logging"
	azure_serverless "github.com/naemono/go-cloud-actions/pkg/serverless/azure"
)

var (
	// AzureCmd is the base azure compute command
	AzureCmd = &cobra.Command{
		Use:              "azure",
		Short:            "Control compute in azure's public clouds",
		Long:             `A cli to interact with compute in Azure's public cloud.`,
		PersistentPreRun: shared_azure.PersistentPreRun,
	}
	computeCreateContainerInstanceCmd = &cobra.Command{
		Use:   "create-container-instance",
		Short: "create compute container instance in azure's public clouds",
		Long:  `A cli to create compute container instance in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("file", cmd.Flags().Lookup("file"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("file") == "" {
				return fmt.Errorf("file name cannot be empty")
			}
			return createContainersGroup()
		},
	}
)

func init() {
	shared_azure.AddAuthFlagsToCommand(AzureCmd)

	computeCreateContainerInstanceCmd.Flags().StringP("file", "f", "", "container yaml file to deploy")

	AzureCmd.AddCommand(computeCreateContainerInstanceCmd)
}

func createContainersGroup() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("creating containers group")
	client, err := azure_serverless.New(azure_serverless.Config{
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var req azure_serverless.CreateContainerRequest
	req, err = readContainersFile(viper.GetString("file"))
	if err != nil {
		return err
	}
	logger.Infof("creating container group with request: %+v", req)
	var cg containerinstance.ContainerGroup
	cg, err = client.CreateContainerGroup(ctx, req)
	if err != nil {
		return err
	}
	logger.Infof("container group '%s' created", *cg.Name)
	return nil
}

func readContainersFile(filename string) (request azure_serverless.CreateContainerRequest, err error) {
	var fileBytes []byte
	fileBytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return request, errors.Wrap(err, "failed to read file")
	}
	temp := map[string]interface{}{}
	err = yaml.Unmarshal(fileBytes, &temp)
	if err != nil {
		return request, errors.Wrap(err, "failed to encode containers yaml file into valid map")
	}
	fileBytes, err = json.Marshal(&temp)
	if err != nil {
		return request, errors.Wrap(err, "failed to marshal map back to json")
	}
	err = json.Unmarshal(fileBytes, &request)
	if err != nil {
		return request, errors.Wrap(err, "failed to encode containers yaml file into valid json")
	}
	return
}
