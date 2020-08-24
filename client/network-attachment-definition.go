package client

import (
	cniresources "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetNetworkAttachmentDefinition(namespace string, networkName string) (*cniresources.NetworkAttachmentDefinition, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	if networkAttachmentDefinition, err := self.Net.K8sCniCncfIoV1().NetworkAttachmentDefinitions(namespace).Get(self.Context, networkName, meta.GetOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if networkAttachmentDefinition.Kind == "" {
			networkAttachmentDefinition = networkAttachmentDefinition.DeepCopy()
			networkAttachmentDefinition.APIVersion = cniresources.SchemeGroupVersion.String()
			networkAttachmentDefinition.Kind = "NetworkAttachmentDefinition"
		}
		return networkAttachmentDefinition, nil
	} else {
		return nil, err
	}
}

func (self *Client) CreateNetworkAttachmentDefinition(network *resources.Network, config string) (*cniresources.NetworkAttachmentDefinition, error) {
	networkAttachmentDefinition := &cniresources.NetworkAttachmentDefinition{
		ObjectMeta: meta.ObjectMeta{
			Name:      network.Name,
			Namespace: network.Namespace,
			OwnerReferences: []meta.OwnerReference{
				*meta.NewControllerRef(network, network.GroupVersionKind()),
			},
		},
		Spec: cniresources.NetworkAttachmentDefinitionSpec{
			Config: config,
		},
	}

	self.Log.Infof("creating network attachment definition: %s/%s\n%s", networkAttachmentDefinition.Namespace, networkAttachmentDefinition.Name, networkAttachmentDefinition.Spec.Config)

	if networkAttachmentDefinition, err := self.Net.K8sCniCncfIoV1().NetworkAttachmentDefinitions(networkAttachmentDefinition.Namespace).Create(self.Context, networkAttachmentDefinition, meta.CreateOptions{}); err == nil {
		return networkAttachmentDefinition, nil
	} else if errors.IsAlreadyExists(err) {
		return self.Net.K8sCniCncfIoV1().NetworkAttachmentDefinitions(network.Namespace).Get(self.Context, network.Name, meta.GetOptions{})
	} else {
		return nil, err
	}
}
