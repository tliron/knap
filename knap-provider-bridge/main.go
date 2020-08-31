package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/heptiolabs/healthcheck"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

var logTo string = ""
var verbose int = 0
var colorize string = "true"

var healthPort uint = 8086

func main() {
	err := terminal.ProcessColorizeFlag(colorize)
	util.FailOnError(err)
	if logTo == "" {
		util.ConfigureLogging(verbose, nil)
	} else {
		util.ConfigureLogging(verbose, &logTo)
	}

	if (len(os.Args) == 1) && os.Args[0] == "provide" {
	} else {
		log.Info("starting health monitor")
		health := healthcheck.NewHandler()
		err = http.ListenAndServe(fmt.Sprintf(":%d", healthPort), health)
		util.FailOnError(err)
	}
}
