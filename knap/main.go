package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/knap/knap/commands"

	_ "github.com/tliron/kutil/logging/simple"
)

func main() {
	commands.Execute()
	atexit.Exit(0)
}
