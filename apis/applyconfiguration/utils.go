// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	knapgithubcomv1alpha1 "github.com/tliron/knap/apis/applyconfiguration/knap.github.com/v1alpha1"
	v1alpha1 "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=knap.github.com, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("Network"):
		return &knapgithubcomv1alpha1.NetworkApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NetworkSpec"):
		return &knapgithubcomv1alpha1.NetworkSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NetworkStatus"):
		return &knapgithubcomv1alpha1.NetworkStatusApplyConfiguration{}

	}
	return nil
}
