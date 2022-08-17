package commands

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var tail int
var follow bool

func init() {
	rootCommand.AddCommand(logsCommand)
	logsCommand.Flags().IntVarP(&tail, "tail", "t", -1, "number of most recent lines to print (<0 means all lines)")
	logsCommand.Flags().BoolVarP(&follow, "follow", "f", false, "keep printing incoming logs")
}

var logsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs of the Knap operator",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: what happens if we follow more than one log?
		readers, err := NewClient().Logs("operator", "operator", tail, follow)
		util.FailOnError(err)
		for _, reader := range readers {
			defer reader.Close()
		}
		for _, reader := range readers {
			io.Copy(os.Stdout, reader)
		}
	},
}
