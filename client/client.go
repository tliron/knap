package client

import (
	"context"

	netpkg "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	"github.com/op/go-logging"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	Config        *rest.Config
	Kubernetes    *kubernetespkg.Clientset
	APIExtensions *apiextensionspkg.Clientset
	Net           *netpkg.Clientset

	Cluster   bool
	Namespace string

	Context context.Context
	Log     *logging.Logger
}
