package azure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	shared_azure "github.com/naemono/go-cloud-actions/cmd/shared/azure"
	auth_azure "github.com/naemono/go-cloud-actions/pkg/auth/azure"
	azure_identity "github.com/naemono/go-cloud-actions/pkg/identity/azure"
	"github.com/naemono/go-cloud-actions/pkg/logging"
)

var (
	AzureCmd = &cobra.Command{
		Use:              "azure",
		Short:            "Control identity in azure's public clouds",
		Long:             `A cli to interact with identity in Azure's public cloud.`,
		PersistentPreRun: shared_azure.PersistentPreRun,
	}
	usersCmd = &cobra.Command{
		Use:   "users",
		Short: "control users in azure's public clouds",
		Long:  `A cli to control users in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
		},
	}
	userAddCmd = &cobra.Command{
		Use:   "add",
		Short: "add users in azure's public clouds",
		Long:  `A cli to add users in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("app-id", cmd.Flags().Lookup("app-id"))
			viper.BindPFlag("display-name", cmd.Flags().Lookup("display-name"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("app-id") == "" {
				return fmt.Errorf("app-id cannot be empty")
			}
			if viper.GetString("display-name") == "" {
				return fmt.Errorf("display-name cannot be empty")
			}
			return createUser()
		},
	}
	applicationsCmd = &cobra.Command{
		Use:   "applications",
		Short: "control applications in azure's public clouds",
		Long:  `A cli to control applicadtions in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
		},
	}
	applicationAddCmd = &cobra.Command{
		Use:   "add",
		Short: "add applications in azure's public clouds",
		Long:  `A cli to add applications in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("multi-tenant", cmd.Flags().Lookup("multi-tenant"))
			viper.BindPFlag("display-name", cmd.Flags().Lookup("display-name"))
			viper.BindPFlag("homepage", cmd.Flags().Lookup("homepage"))
			viper.BindPFlag("identifier-uris", cmd.Flags().Lookup("identifier-uris"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("display-name") == "" {
				return errors.New("display-name cannot be empty")
			}
			return createApplication()
		},
	}
	applicationAddCredentialsCmd = &cobra.Command{
		Use:   "add-credentials",
		Short: "add credentials to an application in azure's public clouds",
		Long:  `A cli to add credentials to an application in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("app-id", cmd.Flags().Lookup("app-id"))
			viper.BindPFlag("display-name", cmd.Flags().Lookup("display-name"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("app-id") == "" {
				return fmt.Errorf("app-id cannot be empty")
			}
			if viper.GetString("display-name") == "" {
				return errors.New("display-name cannot be empty")
			}
			return updateApplicationCredentials()
		},
	}
	rolesCmd = &cobra.Command{
		Use:   "roles",
		Short: "control roles in azure's public clouds",
		Long:  `A cli to control roles in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("resource-group") == "" {
				return fmt.Errorf("resource-group cannot be empty")
			}
			if viper.GetString("vnet-name") == "" {
				return errors.New("vnet-name cannot be empty")
			}
			return nil
		},
	}
	rolesListCmd = &cobra.Command{
		Use:   "list",
		Short: "list roles in azure's public clouds",
		Long:  `A cli to list roles in Azure's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("resource-group", cmd.Flags().Lookup("resource-group"))
			viper.BindPFlag("vnet-name", cmd.Flags().Lookup("vnet-name"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return rolesList()
		},
	}
)

func init() {
	shared_azure.AddAuthFlagsToCommand(AzureCmd)

	applicationAddCmd.Flags().BoolP("multi-tenant", "m", true, "is this app multi-tenant?")
	applicationAddCmd.Flags().StringP("display-name", "d", "", "display name of application")
	applicationAddCmd.Flags().StringP("homepage", "H", "https://microsoft.com", "home page of application")
	applicationAddCmd.Flags().StringSliceP("identifier-uris", "i", []string{}, "list of identifier uris for the application")

	applicationAddCredentialsCmd.Flags().StringP("app-id", "a", "", "application id to add credentials")
	applicationAddCredentialsCmd.Flags().StringP("display-name", "d", "", "display name of application")

	userAddCmd.Flags().StringP("app-id", "a", "", "application id to which to add this user")
	userAddCmd.Flags().StringP("display-name", "d", "", "display name of application")

	rolesListCmd.Flags().StringP("resource-group", "r", "", "resource group to use as scope")
	rolesListCmd.Flags().StringP("vnet-name", "v", "", "vnet name to use as scope")

	AzureCmd.AddCommand(usersCmd)
	AzureCmd.AddCommand(applicationsCmd)
	AzureCmd.AddCommand(rolesCmd)
	usersCmd.AddCommand(userAddCmd)
	applicationsCmd.AddCommand(applicationAddCmd)
	applicationsCmd.AddCommand(applicationAddCredentialsCmd)
	rolesCmd.AddCommand(rolesListCmd)
}

func createUser() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("creating user")
	client := azure_identity.New(azure_identity.Config{
		AuthConfig: auth_azure.AuthConfig{
			SubscriptionID: viper.GetString("subscription-id"),
			ClientID:       viper.GetString("client-id"),
			ClientSecret:   viper.GetString("client-secret"),
			TenantID:       viper.GetString("tenant-id"),
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	appID := viper.GetString("app-id")
	sp, err := client.CreateServicePrincipal(ctx, azure_identity.ApplicationConfig{
		AppID:       &appID,
		DisplayName: viper.GetString("display-name"),
	})
	if err != nil {
		return err
	}
	logger.Infof("service principal created for app '%s'", *sp.DisplayName)
	return nil
}

func createApplication() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("creating application")
	client := azure_identity.New(azure_identity.Config{
		AuthConfig: auth_azure.AuthConfig{
			SubscriptionID: viper.GetString("subscription-id"),
			ClientID:       viper.GetString("client-id"),
			ClientSecret:   viper.GetString("client-secret"),
			TenantID:       viper.GetString("tenant-id"),
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	app, err := client.CreateADApplication(ctx, azure_identity.ApplicationConfig{
		AvailableToOtherTenants: viper.GetBool("multi-tenant"),
		DisplayName:             viper.GetString("display-name"),
		HomePage:                viper.GetString("homepage"),
		IdentifierUris:          viper.GetStringSlice("identifier-uris"),
	})
	if err != nil {
		return err
	}
	logger.Infof("application id %s created", *app.AppID)
	return nil
}

func updateApplicationCredentials() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("updating application credentials")
	client := azure_identity.New(azure_identity.Config{
		AuthConfig: auth_azure.AuthConfig{
			SubscriptionID: viper.GetString("subscription-id"),
			ClientID:       viper.GetString("client-id"),
			ClientSecret:   viper.GetString("client-secret"),
			TenantID:       viper.GetString("tenant-id"),
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	appID := viper.GetString("app-id")
	password, err := client.CreateApplicationCredentials(ctx, azure_identity.ApplicationConfig{
		AppID:       &appID,
		DisplayName: viper.GetString("display-name"),
	})
	if err != nil {
		return err
	}
	logger.Infof("password of %s assigned to application id %s", password, appID)
	return nil
}

func rolesList() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("listing roles")
	client := azure_identity.New(azure_identity.Config{
		AuthConfig: auth_azure.AuthConfig{
			SubscriptionID: viper.GetString("subscription-id"),
			ClientID:       viper.GetString("client-id"),
			ClientSecret:   viper.GetString("client-secret"),
			TenantID:       viper.GetString("tenant-id"),
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	roles, err := client.ListRoleDefinitions(ctx, viper.GetString("resource-group"), viper.GetString("vnet-name"))
	if err != nil {
		return err
	}
	for _, role := range roles {
		logger.Infof("role name: %s, description: %s", *role.Name, *role.Description)
	}
	return nil
}
