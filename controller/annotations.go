package controller

import (
	"fmt"
	"strings"

	cniresources "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ObjectHasNetwork(object *meta.ObjectMeta, name string) bool {
	if networkNames, ok := GetObjectNetworkNames(object); ok {
		for _, networkName := range networkNames {
			if networkName == name {
				return true
			}
		}
	}
	return false
}

func ObjectHasNetworkAttachmentDefinition(object *meta.ObjectMeta, name string) bool {
	if networkAttachmentDefinitionNames, ok := GetObjectNetworkAttachmentDefinition(object); ok {
		for _, networkAttachmentDefinitionName := range networkAttachmentDefinitionNames {
			if networkAttachmentDefinitionName == name {
				return true
			}
		}
	}
	return false
}

func GetObjectNetworkNames(object *meta.ObjectMeta) ([]string, bool) {
	if annotation, ok := object.Annotations[resources.NetworkAnnotation]; ok {
		var networkNames []string
		for _, networkName := range strings.Split(annotation, ",") {
			networkName = strings.TrimSpace(networkName)
			networkNames = append(networkNames, networkName)
		}
		return networkNames, true
	} else {
		return nil, false
	}
}

func GetObjectNetworkAttachmentDefinition(object *meta.ObjectMeta) ([]string, bool) {
	if annotation, ok := object.Annotations[cniresources.NetworkAttachmentAnnot]; ok {
		var networkNames []string
		for _, networkName := range strings.Split(annotation, ",") {
			networkName = strings.TrimSpace(networkName)
			networkNames = append(networkNames, networkName)
		}
		return networkNames, true
	} else {
		return nil, false
	}
}

func AddNetworkAttachmentDefinitionToObject(object *meta.ObjectMeta, name string) {
	if existing, ok := object.Annotations[cniresources.NetworkAttachmentAnnot]; ok {
		object.Annotations[cniresources.NetworkAttachmentAnnot] = fmt.Sprintf("%s,%s", existing, name)
	} else {
		object.Annotations[cniresources.NetworkAttachmentAnnot] = name
	}
}
