package google

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	google_auth "github.com/naemono/go-cloud-actions/pkg/auth/google"
	"github.com/naemono/go-cloud-actions/pkg/logging"
	serverless_google "github.com/naemono/go-cloud-actions/pkg/serverless/google"
	"github.com/naemono/go-cloud-actions/pkg/validate"
)

var (
	// GoogleCmd is the base google compute's command
	GoogleCmd = &cobra.Command{
		Use:   "google",
		Short: "Control compute/gke  in google's public clouds",
		Long:  `A cli to interact with compute/gke in Google's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("google-credentials-file-path", cmd.Flags().Lookup("google-credentials-file-path"))
		},
	}
	createCmd = &cobra.Command{
		Use:   "create-cluster",
		Short: "create gke cluster in google's public clouds",
		Long:  `A cli to create gke clusters in Google's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
				cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			}
			viper.BindPFlag("project-id", cmd.Flags().Lookup("project-id"))
			viper.BindPFlag("network-name", cmd.Flags().Lookup("network-name"))
			viper.BindPFlag("cluster-ipv4-cidr", cmd.Flags().Lookup("cluster-ipv4-cidr"))
			viper.BindPFlag("description", cmd.Flags().Lookup("description"))
			viper.BindPFlag("location", cmd.Flags().Lookup("location"))
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(
				viper.GetViper(),
				[]string{"google-credentials-file-path", "project-id", "network-name", "cluster-ipv4-cidr", "description", "location", "name"}); err != nil {
				return err
			}
			return createCluster()
		},
	}
)

func init() {
	GoogleCmd.PersistentFlags().StringP("google-credentials-file-path", "G", "", "google service account credentials json file")

	createCmd.Flags().StringP("project-id", "p", "", "google project id/name")
	createCmd.Flags().StringP("network-name", "n", "", "google project network name")
	createCmd.Flags().StringP("cluster-ipv4-cidr", "c", "", "ipv4 cidr of cluster")
	createCmd.Flags().StringP("description", "d", "", "description of cluster")
	createCmd.Flags().StringP("location", "L", "", "location in which to create a cluster")
	createCmd.Flags().StringP("name", "N", "", "name of the cluster to create")

	GoogleCmd.AddCommand(createCmd)
}

func createCluster() error {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	logger.Infof("creating cluster")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := serverless_google.New(serverless_google.Config{
		AuthConfig: google_auth.AuthConfig{
			CredentialsFilePath: viper.GetString("google-credentials-file-path"),
		},
		Logger: logger,
	})
	if err != nil {
		return err
	}
	return client.CreateCluster(ctx, serverless_google.CreateClusterRequest{
		ClusterCommon: serverless_google.ClusterCommon{
			ProjectID:   viper.GetString("project-id"),
			NetworkName: viper.GetString("network-name"),
		},
		ClusterIpv4Cidr: viper.GetString("cluster-ipv4-cidr"),
		Description:     viper.GetString("description"),
		Location:        viper.GetString("location"),
		Name:            viper.GetString("name"),
	})
}
