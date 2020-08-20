// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/tliron/knap/apis/clientset/versioned/typed/knap.github.com/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeKnapV1alpha1 struct {
	*testing.Fake
}

func (c *FakeKnapV1alpha1) Networks(namespace string) v1alpha1.NetworkInterface {
	return &FakeNetworks{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKnapV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
