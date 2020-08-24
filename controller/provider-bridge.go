package controller

import (
	"errors"
	"fmt"
	"sync"

	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
)

var bridgePrefixPool = []string{
	"192.168.2",
	"192.168.3",
	"192.168.4",
	"192.168.5",
	"192.168.6",
	"192.168.7",
}
var bridgePrefixLock sync.Mutex

func getBridgePrefix() (string, error) {
	bridgePrefixLock.Lock()
	defer bridgePrefixLock.Unlock()

	if len(bridgePrefixPool) == 0 {
		return "", errors.New("bridge prefix pool is empty")
	}

	bridgePrefix := bridgePrefixPool[0]
	bridgePrefixPool = bridgePrefixPool[1:]
	return bridgePrefix, nil
}

func (self *Controller) createBridgeCniConfig(network *resources.Network) (string, error) {
	if bridgePrefix, err := getBridgePrefix(); err == nil {
		return fmt.Sprintf(`{
  "cniVersion": "0.3.1",
  "type": "bridge",
  "name": "%s",
  "bridge": "%s",
  "isDefaultGateway": true,
  "ipMasq": true,
  "promiscMode": true,
  "ipam": {
    "type": "host-local",
    "subnet": "%s.0/24",
    "rangeStart": "%s.2",
    "rangeEnd": "%s.254",
    "routes": [
      { "dst": "0.0.0.0/0" }
    ],
    "gateway": "%s.1"
  }
}`, network.Name, network.Name, bridgePrefix, bridgePrefix, bridgePrefix, bridgePrefix), nil
	} else {
		return "", err
	}
}
