package cmd

import (
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/naemono/go-cloud-actions/cmd/compute"
	"github.com/naemono/go-cloud-actions/cmd/identity"
	"github.com/naemono/go-cloud-actions/cmd/network"
	"github.com/naemono/go-cloud-actions/cmd/peering"
	"github.com/naemono/go-cloud-actions/cmd/resources"
	"github.com/naemono/go-cloud-actions/pkg/logging"
)

var (
	// CloudCmd is the root cloud command
	CloudCmd = &cobra.Command{
		Use:     "cloud",
		Version: version,
		Short:   "Control peering, and user credentials in public clouds",
		Long:    `A cli to interact with users, permissions, and peering in public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	version string
)

func init() {
	CloudCmd.PersistentFlags().StringP("loglevel", "l", "info", "logging level")
	viper.BindPFlag("loglevel", CloudCmd.PersistentFlags().Lookup("loglevel"))
	CloudCmd.AddCommand(compute.RootCmd)
	CloudCmd.AddCommand(peering.RootCmd)
	CloudCmd.AddCommand(identity.RootCmd)
	CloudCmd.AddCommand(resources.RootCmd)
	CloudCmd.AddCommand(network.RootCmd)
}

// Run will run the main command
func Run() {
	logging.GetLogger(viper.GetString("loglevel")).Infof("running cloud version: %s", version)
	if err := CloudCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("failure running cloud command")
	}
}
