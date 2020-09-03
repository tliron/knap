package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var cluster bool
var registry string
var wait bool

func init() {
	rootCommand.AddCommand(installCommand)
	installCommand.Flags().BoolVarP(&cluster, "cluster", "c", false, "cluster mode")
	installCommand.Flags().StringVarP(&registry, "registry", "g", "docker.io", "registry URL (use special value \"internal\" to discover internally deployed registry)")
	installCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for installation to succeed")
}

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Install Knap",
	Run: func(cmd *cobra.Command, args []string) {
		err := NewClient().Install(registry, wait)
		util.FailOnError(err)
	},
}
