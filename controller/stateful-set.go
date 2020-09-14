package controller

import (
	cniresources "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Controller) processStatefulSets(network *resources.Network, networkAttachmentDefinition *cniresources.NetworkAttachmentDefinition) error {
	if statefulSets, err := self.Kubernetes.AppsV1().StatefulSets(networkAttachmentDefinition.Namespace).List(self.Context, meta.ListOptions{}); err == nil {
		for _, statefulSet := range statefulSets.Items {
			object := &statefulSet.Spec.Template.ObjectMeta

			if ObjectHasNetwork(object, network.Name) {
				self.Log.Infof("processing stateful set %s/%s for network %q", statefulSet.Namespace, statefulSet.Name, network.Name)

				if !ObjectHasNetworkAttachmentDefinition(object, networkAttachmentDefinition.Name) {
					statefulSet_ := statefulSet.DeepCopy()
					object = &statefulSet_.Spec.Template.ObjectMeta
					AddNetworkAttachmentDefinitionToObject(object, networkAttachmentDefinition.Name)
					if _, err := self.Kubernetes.AppsV1().StatefulSets(networkAttachmentDefinition.Namespace).Update(self.Context, statefulSet_, meta.UpdateOptions{}); err == nil {
						self.Log.Infof("attached stateful set %s/%s to network attachment definition %q", statefulSet.Namespace, statefulSet.Name, networkAttachmentDefinition.Name)
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
