package network

import (
	"github.com/spf13/cobra"

	"github.com/naemono/go-cloud-actions/cmd/network/aws"
	"github.com/naemono/go-cloud-actions/cmd/network/azure"
)

var (
	// RootCmd is the base network command for all public clouds
	RootCmd = &cobra.Command{
		Use:   "network",
		Short: "Control networks in public clouds",
		Long:  `A cli to interact with networks in AWS, Azure, and GCP public clouds.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(azure.AzureCmd)
	RootCmd.AddCommand(aws.AWSCmd)
}
