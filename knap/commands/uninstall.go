package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(uninstallCommand)
	uninstallCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for uninstallation to succeed")
}

var uninstallCommand = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Knap",
	Run: func(cmd *cobra.Command, args []string) {
		NewClient().Uninstall(wait)
	},
}
