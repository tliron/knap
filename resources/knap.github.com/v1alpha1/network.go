package v1alpha1

import (
	"fmt"

	group "github.com/tliron/knap/resources/knap.github.com"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var NetworkGVK = SchemeGroupVersion.WithKind(NetworkKind)

const (
	NetworkKind     = "Network"
	NetworkListKind = "NetworkList"

	NetworkSingular  = "network"
	NetworkPlural    = "networks"
	NetworkShortName = "nw"
)

//
// Network
//

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Network struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkSpec   `json:"spec"`
	Status NetworkStatus `json:"status"`
}

type NetworkSpec struct {
	Name string `json:"name"`
}

type NetworkStatus struct {
	Name string `json:"name"`
}

//
// NetworkList
//

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NetworkList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata"`

	Items []Network `json:"items"`
}

//
// NetworkCustomResourceDefinition
//

// See: assets/custom-resource-definitions.yaml

var NetworkResourcesName = fmt.Sprintf("%s.%s", NetworkPlural, group.GroupName)

var NetworkCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: NetworkResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: NetworkSingular,
			Plural:   NetworkPlural,
			Kind:     NetworkKind,
			ListKind: NetworkListKind,
			ShortNames: []string{
				NetworkShortName,
			},
			Categories: []string{
				"all", // will appear in "kubectl get all"
			},
		},
		Scope: apiextensions.NamespaceScoped,
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    Version,
				Served:  true,
				Storage: true, // one and only one version must be marked with storage=true
				Subresources: &apiextensions.CustomResourceSubresources{ // requires CustomResourceSubresources feature gate enabled
					Status: &apiextensions.CustomResourceSubresourceStatus{},
				},
				Schema: &apiextensions.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
						Type:     "object",
						Required: []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type:     "object",
								Required: []string{"name"},
								Properties: map[string]apiextensions.JSONSchemaProps{
									"name": {
										Type: "string",
									},
								},
							},
							"status": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"name": {
										Type: "string",
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
