package commands

import (
	"github.com/tliron/knap/version"
)

func init() {
	rootCommand.AddCommand(version.NewCommand(toolName))
}
