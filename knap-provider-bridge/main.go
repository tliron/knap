package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/heptiolabs/healthcheck"
	"github.com/tliron/commonlog"
	"github.com/tliron/go-ard"
	"github.com/tliron/knap/knap-provider-bridge/server"
	"github.com/tliron/knap/provider"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"k8s.io/klog/v2"

	_ "github.com/tliron/commonlog/simple"
)

const SOCKET_NAME = "/tmp/knap-provider-bridge.sock"
const STATE_FILENAME = "/tmp/knap-provider-bridge.state"
const HEALTH_PORT uint = 8086

var logTo string = ""
var verbose int = 1
var colorize string = "true"

func main() {
	cleanup, err := terminal.ProcessColorizeFlag(colorize)
	util.FailOnError(err)
	if cleanup != nil {
		util.OnExitError(cleanup)
	}
	if logTo == "" {
		commonlog.Configure(verbose, nil)
	} else {
		commonlog.Configure(verbose, &logTo)
	}
	if writer := commonlog.GetWriter(); writer != nil {
		klog.SetOutput(writer)
	}

	if (len(os.Args) == 3) && os.Args[1] == "provide" {
		Client(os.Args[2])
	} else {
		Server()
	}
}

// Output a CNI configuration to stdout
func Client(name string) {
	hints := GetHints()

	client, err := provider.NewClient(SOCKET_NAME, log)
	util.FailOnError(err)
	defer client.Release()

	cni, err := client.CreateCniConfig(name, hints)
	util.FailOnError(err)
	fmt.Fprintln(os.Stdout, cni)
}

func GetHints() map[string]string {
	bytes, err := io.ReadAll(os.Stdin)
	util.FailOnError(err)
	if len(bytes) == 0 {
		return nil
	}

	hints, _, err := ard.DecodeYAML(util.BytesToString(bytes), false)
	util.FailOnError(err)

	var hints_ map[string]string
	if hints__, ok := hints.(ard.Map); ok {
		hints_ = make(map[string]string)
		for key, value := range hints__ {
			if key_, ok := key.(string); ok {
				if value_, ok := value.(string); ok {
					hints_[key_] = value_
				} else {
					util.Fail("malformed hints in stdin")
				}
			} else {
				util.Fail("malformed hints in stdin")
			}
		}
	}

	return hints_
}

// Run infinitely and provide services for the client
func Server() {
	go func() {
		log.Infof("starting health monitor on port %d", HEALTH_PORT)
		health := healthcheck.NewHandler()
		err := http.ListenAndServe(fmt.Sprintf(":%d", HEALTH_PORT), health)
		util.FailOnError(err)
	}()

	err := server.NewServer(STATE_FILENAME, SOCKET_NAME, log).Start()
	util.FailOnError(err)
}
