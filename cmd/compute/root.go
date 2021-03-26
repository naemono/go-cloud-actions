package compute

import (
	"github.com/spf13/cobra"

	"github.com/naemono/go-cloud-actions/cmd/compute/azure"
)

var (
	RootCmd = &cobra.Command{
		Use:   "compute",
		Short: "Control compute in public clouds",
		Long:  `A cli to interact with compute in AWS, Azure, and GCP public clouds.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(azure.AzureCmd)
}
