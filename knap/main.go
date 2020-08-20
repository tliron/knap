package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/knap/knap/commands"
)

func main() {
	commands.Execute()
	atexit.Exit(0)
}
