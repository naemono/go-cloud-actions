package shared

import "github.com/spf13/cobra"

// RunParentsPersistentPreRun will be used to ensure that the parent's
// persistent pre-run is run for every cobra command
func RunParentsPersistentPreRun(cmd *cobra.Command, args []string) {
	if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
		cmd.Parent().PersistentPreRun(cmd.Parent(), args)
	}
}
