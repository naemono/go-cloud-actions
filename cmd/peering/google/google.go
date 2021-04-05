package google

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	google_auth "github.com/naemono/go-cloud-actions/pkg/auth/google"
	"github.com/naemono/go-cloud-actions/pkg/logging"
	peering_google "github.com/naemono/go-cloud-actions/pkg/peering/google"
	"github.com/naemono/go-cloud-actions/pkg/validate"
)

var (
	// GoogleCmd is the base google peering's command
	GoogleCmd = &cobra.Command{
		Use:   "google",
		Short: "Control peering of VPCs in google's public clouds",
		Long:  `A cli to interact with peering of VPCs in Google's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("google-credentials-file-path", cmd.Flags().Lookup("google-credentials-file-path"))
		},
	}
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create peers of VPCs in google's public clouds",
		Long:  `A cli to create peers of VPCs in Google's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("project-id", cmd.Flags().Lookup("project-id"))
			viper.BindPFlag("network-name", cmd.Flags().Lookup("network-name"))
			viper.BindPFlag("peering-name", cmd.Flags().Lookup("peering-name"))
			viper.BindPFlag("remote-project-name", cmd.Flags().Lookup("remote-project-name"))
			viper.BindPFlag("remote-network-name", cmd.Flags().Lookup("remote-network-name"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(
				viper.GetViper(),
				[]string{
					"google-credentials-file-path", "project-id", "network-name", "peering-name",
					"remote-project-name", "remote-network-name"}); err != nil {
				return err
			}
			return createPeering()
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list peers of VPCs in google's public clouds",
		Long:  `A cli to list peers of VPCs in Google's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("project-id", cmd.Flags().Lookup("project-id"))
			viper.BindPFlag("network-name", cmd.Flags().Lookup("network-name"))
			viper.BindPFlag("region", cmd.Flags().Lookup("region"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("google-credentials-file-path") == "" {
				return fmt.Errorf("google-credentials-file-path cannot be empty")
			}
			if viper.GetString("project-id") == "" {
				return fmt.Errorf("project-id cannot be empty")
			}
			if viper.GetString("network-name") == "" {
				return fmt.Errorf("network-name cannot be empty")
			}
			if viper.GetString("region") == "" {
				return fmt.Errorf("region cannot be empty")
			}
			return listPeerings()
		},
	}
)

func init() {
	GoogleCmd.PersistentFlags().StringP("google-credentials-file-path", "G", "", "google service account credentials json file")

	createCmd.Flags().StringP("project-id", "p", "", "google project id/name")
	createCmd.Flags().StringP("network-name", "n", "", "google project network name")
	createCmd.Flags().StringP("peering-name", "P", "", "peering name to create")
	createCmd.Flags().StringP("remote-project-name", "r", "", "google remote project name to peer with")
	createCmd.Flags().StringP("remote-network-name", "R", "", "google project network name to peer with")

	listCmd.Flags().StringP("project-id", "p", "", "google project id/name")
	listCmd.Flags().StringP("network-name", "n", "", "google project network name")
	listCmd.Flags().StringP("region", "r", "", "google project network region")

	GoogleCmd.AddCommand(createCmd)
	GoogleCmd.AddCommand(listCmd)
}

func createPeering() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("creating peering")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := peering_google.New(peering_google.Config{
		AuthConfig: google_auth.AuthConfig{
			CredentialsFilePath: viper.GetString("google-credentials-file-path"),
		},
		Logger: logger,
	})
	if err != nil {
		return err
	}
	return client.CreatePeering(ctx, peering_google.CreatePeeringRequest{
		PeeringCommon: peering_google.PeeringCommon{
			PeeringName: viper.GetString("peering-name"),
			ProjectID:   viper.GetString("project-id"),
			NetworkName: viper.GetString("network-name"),
		},
		RemoteNetworkName: viper.GetString("remote-network-name"),
		RemoteProjectName: viper.GetString("remote-project-name"),
	})
}

func listPeerings() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Info("listing peerings")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := peering_google.New(peering_google.Config{
		AuthConfig: google_auth.AuthConfig{
			CredentialsFilePath: viper.GetString("google-credentials-file-path"),
		},
		Logger: logger,
	})
	if err != nil {
		return err
	}
	return client.ListPeerings(ctx, peering_google.ListPeeringRequest{
		PeeringCommon: peering_google.PeeringCommon{
			PeeringName: viper.GetString("peering-name"),
			ProjectID:   viper.GetString("project-id"),
			NetworkName: viper.GetString("network-name"),
		},
		Region: viper.GetString("region"),
	})
}
