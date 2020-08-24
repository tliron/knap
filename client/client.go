package client

import (
	"context"

	netpkg "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	"github.com/op/go-logging"
	knappkg "github.com/tliron/knap/apis/clientset/versioned"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	Config        *rest.Config
	Kubernetes    kubernetespkg.Interface
	APIExtensions apiextensionspkg.Interface
	Net           netpkg.Interface
	Knap          knappkg.Interface

	Cluster           bool
	Namespace         string
	NamePrefix        string
	PartOf            string
	ManagedBy         string
	OperatorImageName string

	Context context.Context
	Log     *logging.Logger
}
