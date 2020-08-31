package commands

import (
	contextpkg "context"

	netpkg "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	knappkg "github.com/tliron/knap/apis/clientset/versioned"
	"github.com/tliron/knap/client"
	"github.com/tliron/knap/controller"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
)

func NewClient() *client.Client {
	config, err := kubernetesutil.NewConfigFromFlags(masterUrl, kubeconfigPath, context, log)
	util.FailOnError(err)

	kubernetes, err := kubernetespkg.NewForConfig(config)
	util.FailOnError(err)

	apiExtensions, err := apiextensionspkg.NewForConfig(config)
	util.FailOnError(err)

	net, err := netpkg.NewForConfig(config)
	util.FailOnError(err)

	knap, err := knappkg.NewForConfig(config)
	util.FailOnError(err)

	namespace_ := namespace
	if cluster {
		namespace_ = ""
	} else if namespace_ == "" {
		if namespace__, ok := kubernetesutil.GetConfiguredNamespace(kubeconfigPath, context); ok {
			namespace_ = namespace__
		}
		if namespace_ == "" {
			util.Fail("could not discover namespace and \"--namespace\" not provided")
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
