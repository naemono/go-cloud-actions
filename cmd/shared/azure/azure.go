package azure

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PersistentPreRun is a shared persistent pre-run for azure commands
func PersistentPreRun(cmd *cobra.Command, args []string) {
	if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
		cmd.Parent().PersistentPreRun(cmd.Parent(), args)
	}
	viper.BindPFlag("client-id", cmd.Flags().Lookup("client-id"))
	viper.BindPFlag("client-secret", cmd.Flags().Lookup("client-secret"))
	viper.BindPFlag("subscription-id", cmd.Flags().Lookup("subscription-id"))
	viper.BindPFlag("tenant-id", cmd.Flags().Lookup("tenant-id"))
}

func AddAuthFlagsToCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("client-id", "c", "", "azure client id")
	cmd.PersistentFlags().StringP("client-secret", "S", "", "azure client secret")
	cmd.PersistentFlags().StringP("subscription-id", "s", "", "azure subscription id")
	cmd.PersistentFlags().StringP("tenant-id", "t", "", "azure tenant id")
}
