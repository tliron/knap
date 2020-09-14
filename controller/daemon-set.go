package controller

import (
	cniresources "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Controller) processDaemonSets(network *resources.Network, networkAttachmentDefinition *cniresources.NetworkAttachmentDefinition) error {
	if daemonSets, err := self.Kubernetes.AppsV1().DaemonSets(networkAttachmentDefinition.Namespace).List(self.Context, meta.ListOptions{}); err == nil {
		for _, daemonSet := range daemonSets.Items {
			object := &daemonSet.Spec.Template.ObjectMeta

			if ObjectHasNetwork(object, network.Name) {
				self.Log.Infof("processing daemon set %s/%s for network %q", daemonSet.Namespace, daemonSet.Name, network.Name)

				if !ObjectHasNetworkAttachmentDefinition(object, networkAttachmentDefinition.Name) {
					daemonSet_ := daemonSet.DeepCopy()
					object = &daemonSet_.Spec.Template.ObjectMeta
					AddNetworkAttachmentDefinitionToObject(object, networkAttachmentDefinition.Name)
					if _, err := self.Kubernetes.AppsV1().DaemonSets(networkAttachmentDefinition.Namespace).Update(self.Context, daemonSet_, meta.UpdateOptions{}); err == nil {
						self.Log.Infof("attached daemon set %s/%s to network attachment definition %q", daemonSet.Namespace, daemonSet.Name, networkAttachmentDefinition.Name)
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
