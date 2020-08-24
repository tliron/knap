package client

import (
	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetNetwork(namespace string, networkName string) (*resources.Network, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	if network, err := self.Knap.KnapV1alpha1().Networks(namespace).Get(self.Context, networkName, meta.GetOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if network.Kind == "" {
			network = network.DeepCopy()
			network.APIVersion, network.Kind = resources.NetworkGVK.ToAPIVersionAndKind()
		}
		return network, nil
	} else {
		return nil, err
	}
}

func (self *Client) ListNetworks() (*resources.NetworkList, error) {
	// TODO: all networks in cluster mode
	return self.Knap.KnapV1alpha1().Networks(self.Namespace).List(self.Context, meta.ListOptions{})
}

func (self *Client) DeleteNetwork(namespace string, serviceName string) error {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	return self.Knap.KnapV1alpha1().Networks(namespace).Delete(self.Context, serviceName, meta.DeleteOptions{})
}
