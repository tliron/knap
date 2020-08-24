package controller

import (
	"fmt"

	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
)

func (self *Controller) createCniConfig(network *resources.Network) (string, error) {
	switch network.Spec.Provider {
	case "bridge":
		return self.createBridgeCniConfig(network)

	default:
		return "", fmt.Errorf("unsupported network provider: %s", network.Spec.Provider)
	}
}
