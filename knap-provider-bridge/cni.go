package main

import (
	"fmt"
)

func createBridgeCniConfig(name string) (string, error) {
	if bridgePrefix, err := GetBridgePrefix(); err == nil {
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
}`, name, name, bridgePrefix, bridgePrefix, bridgePrefix, bridgePrefix), nil
	} else {
		return "", err
	}
}
