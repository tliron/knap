package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/kutil/util"

	_ "github.com/tliron/kutil/logging/simple"
)

func main() {
	err := command.Execute()
	util.FailOnError(err)
	atexit.Exit(0)
}
