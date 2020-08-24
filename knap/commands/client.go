package commands

import (
	contextpkg "context"

	netpkg "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	knappkg "github.com/tliron/knap/apis/clientset/versioned"
	"github.com/tliron/knap/client"
	"github.com/tliron/knap/controller"
	puccinicommon "github.com/tliron/puccini/common"
	turandotcommon "github.com/tliron/turandot/common"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
)

func NewClient() *client.Client {
	config, err := turandotcommon.NewConfigFromFlags(masterUrl, kubeconfigPath, context, log)
	puccinicommon.FailOnError(err)

	kubernetes, err := kubernetespkg.NewForConfig(config)
	puccinicommon.FailOnError(err)

	apiExtensions, err := apiextensionspkg.NewForConfig(config)
	puccinicommon.FailOnError(err)

	net, err := netpkg.NewForConfig(config)
	puccinicommon.FailOnError(err)

	knap, err := knappkg.NewForConfig(config)
	puccinicommon.FailOnError(err)

	namespace_ := namespace
	if cluster {
		namespace_ = ""
	} else if namespace_ == "" {
		if namespace__, ok := turandotcommon.GetConfiguredNamespace(kubeconfigPath, context); ok {
			namespace_ = namespace__
		}
		if namespace_ == "" {
			puccinicommon.Fail("could not discover namespace and \"--namespace\" not provided")
		}
	}

	return &client.Client{
		Config:            config,
		Kubernetes:        kubernetes,
		APIExtensions:     apiExtensions,
		Net:               net,
		Knap:              knap,
		Cluster:           cluster,
		Namespace:         namespace_,
		NamePrefix:        controller.NamePrefix,
		PartOf:            controller.PartOf,
		ManagedBy:         controller.ManagedBy,
		OperatorImageName: controller.OperatorImageName,
		Context:           contextpkg.TODO(),
		Log:               log,
	}
}
