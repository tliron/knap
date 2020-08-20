package main

import (
	"fmt"
	"net/http"

	"github.com/heptiolabs/healthcheck"
	"github.com/tebeka/atexit"
	knappkg "github.com/tliron/knap/apis/clientset/versioned"
	controllerpkg "github.com/tliron/knap/controller"
	versionpkg "github.com/tliron/knap/version"
	puccinicommon "github.com/tliron/puccini/common"
	"github.com/tliron/turandot/common"
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
	puccinicommon.FailOnError(err)

	if cluster {
		namespace = ""
	} else if namespace == "" {
		if namespace_, ok := common.GetConfiguredNamespace(kubeconfigPath, context); ok {
			namespace = namespace_
		}
		if namespace == "" {
			namespace = common.GetServiceAccountNamespace()
		}
		if namespace == "" {
			log.Fatal("could not discover namespace and namespace not provided")
		}
	}

	// Clients

	kubernetesClient, err := kubernetespkg.NewForConfig(config)
	puccinicommon.FailOnError(err)

	apiExtensionsClient, err := apiextensionspkg.NewForConfig(config)
	puccinicommon.FailOnError(err)

	knapClient, err := knappkg.NewForConfig(config)
	puccinicommon.FailOnError(err)

	// Controller

	controller := controllerpkg.NewController(
		toolName,
		cluster,
		namespace,
		kubernetesClient,
		apiExtensionsClient,
		knapClient,
		config,
		resyncPeriod,
		common.SetupSignalHandler(),
	)

	// Run

	err = controller.Run(concurrency, func() {
		log.Info("starting health monitor")
		health := healthcheck.NewHandler()
		err := http.ListenAndServe(fmt.Sprintf(":%d", healthPort), health)
		puccinicommon.FailOnError(err)
	})
	puccinicommon.FailOnError(err)
}
