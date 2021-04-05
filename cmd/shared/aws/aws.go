package aws

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PersistentPreRun is a shared persistent pre-run for aws commands
func PersistentPreRun(cmd *cobra.Command, args []string) {
	if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
		cmd.Parent().PersistentPreRun(cmd.Parent(), args)
	}
	viper.BindPFlag("profile", cmd.Flags().Lookup("profile"))
	viper.BindPFlag("region", cmd.Flags().Lookup("region"))
}

// AddAuthFlagsToCommand is a shared command to add the aws auth components to any aws cobra command
func AddAuthFlagsToCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("profile", "p", "", "aws profile to use")
	cmd.PersistentFlags().StringP("region", "r", "", "aws region")
}
