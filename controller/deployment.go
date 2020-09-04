package controller

import (
	cniresources "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Controller) processDeployments(network *resources.Network, networkAttachmentDefinition *cniresources.NetworkAttachmentDefinition) error {
	if deployments, err := self.Kubernetes.AppsV1().Deployments(networkAttachmentDefinition.Namespace).List(self.Context, meta.ListOptions{}); err == nil {
		for _, deployment := range deployments.Items {
			object := &deployment.Spec.Template.ObjectMeta

			if ObjectHasNetwork(object, network.Name) {
				self.Log.Infof("processing deployment %s/%s for network %q", deployment.Namespace, deployment.Name, network.Name)

				if !ObjectHasNetworkAttachmentDefinition(object, networkAttachmentDefinition.Name) {
					deployment_ := deployment.DeepCopy()
					object = &deployment_.Spec.Template.ObjectMeta
					AddNetworkAttachmentDefinitionToObject(object, networkAttachmentDefinition.Name)
					if _, err := self.Kubernetes.AppsV1().Deployments(networkAttachmentDefinition.Namespace).Update(self.Context, deployment_, meta.UpdateOptions{}); err == nil {
						self.Log.Infof("attached deployment %s/%s to network attachment definition %q", deployment.Namespace, deployment.Name, networkAttachmentDefinition.Name)
					} else {
						return err
					}
				}
			}
		}
		return nil
	} else {
		return err
	}
}
