package main

import (
	"fmt"
	"net/http"

	"github.com/heptiolabs/healthcheck"
	netpkg "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	"github.com/tebeka/atexit"
	knappkg "github.com/tliron/knap/apis/clientset/versioned"
	controllerpkg "github.com/tliron/knap/controller"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	versionpkg "github.com/tliron/kutil/version"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Load all auth plugins:
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func Controller() {
	if version {
		versionpkg.Print()
		atexit.Exit(0)
		return
	}

	log.Noticef("%s version=%s revision=%s", toolName, versionpkg.GitVersion, versionpkg.GitRevision)

	// Config

	config, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	util.FailOnError(err)

	if cluster {
		namespace = ""
	} else if namespace == "" {
		if namespace_, ok := kubernetesutil.GetConfiguredNamespace(kubeconfigPath, context); ok {
			namespace = namespace_
		}
		if namespace == "" {
			namespace = kubernetesutil.GetServiceAccountNamespace()
		}
		if namespace == "" {
			util.Fail("could not discover namespace and namespace not provided")
		}
	}

	// Clients

	kubernetesClient, err := kubernetespkg.NewForConfig(config)
	util.FailOnError(err)

	apiExtensionsClient, err := apiextensionspkg.NewForConfig(config)
	util.FailOnError(err)

	netClient, err := netpkg.NewForConfig(config)
	util.FailOnError(err)

	knapClient, err := knappkg.NewForConfig(config)
	util.FailOnError(err)

	// Controller

	controller := controllerpkg.NewController(
		toolName,
		cluster,
		namespace,
		kubernetesClient,
		apiExtensionsClient,
		netClient,
		knapClient,
		config,
		resyncPeriod,
		util.SetupSignalHandler(),
	)

	// Run

	err = controller.Run(concurrency, func() {
		log.Info("starting health monitor")
		health := healthcheck.NewHandler()
		err := http.ListenAndServe(fmt.Sprintf(":%d", healthPort), health)
		util.FailOnError(err)
	})
	util.FailOnError(err)
}
