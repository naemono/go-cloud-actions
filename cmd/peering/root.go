package peering

import (
	"github.com/spf13/cobra"

	"github.com/naemono/go-cloud-actions/cmd/peering/azure"
)

var (
	RootCmd = &cobra.Command{
		Use:   "peering",
		Short: "Control peering of VPCs/VNets in public clouds",
		Long:  `A cli to interact with peering of VPCs and VNets in AWS, Azure, and GCP public clouds.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(azure.AzureCmd)
}
