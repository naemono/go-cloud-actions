package identity

import (
	"github.com/spf13/cobra"

	"github.com/naemono/go-cloud-actions/cmd/identity/azure"
)

var (
	// RootCmd is the base identity command for all clouds
	RootCmd = &cobra.Command{
		Use:   "identity",
		Short: "Control identity (users and permissions) in public clouds",
		Long:  `A cli to interact with identity (users/permissions) in AWS, Azure, and GCP public clouds.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(azure.AzureCmd)
}
