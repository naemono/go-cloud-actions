package resources

import (
	"github.com/spf13/cobra"

	"github.com/naemono/go-cloud-actions/cmd/resources/azure"
)

var (
	RootCmd = &cobra.Command{
		Use:   "resources",
		Short: "Control resources in public clouds",
		Long:  `A cli to interact with resources (resource groups, etc) in AWS, Azure, and GCP public clouds.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(azure.AzureCmd)
}
