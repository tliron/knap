package controller

import (
	cniresources "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (self *Controller) processNetwork(network *resources.Network) (bool, error) {
	var networkAttachmentDefinition *cniresources.NetworkAttachmentDefinition

	var err error
	if networkAttachmentDefinition, err = self.Client.GetNetworkAttachmentDefinition(network.Namespace, network.Name); err != nil {
		if errors.IsNotFound(err) {
			if cniConfig, err := self.createCniConfig(network); err == nil {
				if networkAttachmentDefinition, err = self.Client.CreateNetworkAttachmentDefinition(network, cniConfig); err != nil {
					return false, err
				}
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	}

	if err := self.processDeployments(network, networkAttachmentDefinition); err == nil {
		return true, nil
	} else {
		return false, err
	}
}
