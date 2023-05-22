package controller

import (
	contextpkg "context"
	"time"

	netpkg "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	"github.com/tliron/commonlog"
	knapclientset "github.com/tliron/knap/apis/clientset/versioned"
	knapinformers "github.com/tliron/knap/apis/informers/externalversions"
	knaplisters "github.com/tliron/knap/apis/listers/knap.github.com/v1alpha1"
	clientpkg "github.com/tliron/knap/client"
	knapresources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

//
// Controller
//

type Controller struct {
	Config      *restpkg.Config
	Kubernetes  kubernetes.Interface
	Knap        knapclientset.Interface
	REST        restpkg.Interface
	Client      *clientpkg.Client
	StopChannel <-chan struct{}

	Processors *kubernetesutil.Processors
	Events     record.EventRecorder

	KubernetesInformerFactory informers.SharedInformerFactory
	KnapInformerFactory       knapinformers.SharedInformerFactory

	Networks knaplisters.NetworkLister

	Context contextpkg.Context
	Log     commonlog.Logger
}

func NewController(toolName string, cluster bool, namespace string, kubernetes kubernetes.Interface, apiExtensions apiextensionspkg.Interface, net netpkg.Interface, knap knapclientset.Interface, config *restpkg.Config, informerResyncPeriod time.Duration, stopChannel <-chan struct{}) *Controller {
	context := contextpkg.TODO()

	if cluster {
		namespace = ""
	}

	log := commonlog.GetLoggerf("%s.controller", toolName)

	self := Controller{
		Config:      config,
		Kubernetes:  kubernetes,
		Knap:        knap,
		REST:        kubernetes.CoreV1().RESTClient(),
		StopChannel: stopChannel,
		Processors:  kubernetesutil.NewProcessors(toolName),
		Events:      kubernetesutil.CreateEventRecorder(kubernetes, "Knap", log),
		Context:     context,
		Log:         log,
	}

	self.Client = &clientpkg.Client{
		Config:            config,
		Kubernetes:        kubernetes,
		APIExtensions:     apiExtensions,
		Net:               net,
		Knap:              knap,
		Cluster:           cluster,
		Namespace:         namespace,
		NamePrefix:        NamePrefix,
		PartOf:            PartOf,
		ManagedBy:         ManagedBy,
		OperatorImageName: OperatorImageName,
		Context:           contextpkg.TODO(),
		Log:               log,
	}

	if cluster {
		self.KubernetesInformerFactory = informers.NewSharedInformerFactory(kubernetes, informerResyncPeriod)
		self.KnapInformerFactory = knapinformers.NewSharedInformerFactory(knap, informerResyncPeriod)
	} else {
		self.KubernetesInformerFactory = informers.NewSharedInformerFactoryWithOptions(kubernetes, informerResyncPeriod, informers.WithNamespace(namespace))
		self.KnapInformerFactory = knapinformers.NewSharedInformerFactoryWithOptions(knap, informerResyncPeriod, knapinformers.WithNamespace(namespace))
	}

	// Informers
	networkInformer := self.KnapInformerFactory.Knap().V1alpha1().Networks()

	// Listers
	self.Networks = networkInformer.Lister()

	// Processors

	processorPeriod := 5 * time.Second

	self.Processors.Add(knapresources.NetworkGVK, kubernetesutil.NewProcessor(
		toolName,
		"networks",
		networkInformer.Informer(),
		processorPeriod,
		func(name string, namespace string) (interface{}, error) {
			return self.Client.GetNetwork(namespace, name)
		},
		func(object interface{}) (bool, error) {
			return self.processNetwork(object.(*knapresources.Network))
		},
	))

	return &self
}

func (self *Controller) Run(concurrency uint, startup func()) error {
	defer utilruntime.HandleCrash()

	self.Log.Info("starting informer factories")
	self.KubernetesInformerFactory.Start(self.StopChannel)
	self.KnapInformerFactory.Start(self.StopChannel)

	self.Log.Info("waiting for processor informer caches to sync")
	utilruntime.HandleError(self.Processors.WaitForCacheSync(self.StopChannel))

	self.Log.Infof("starting processors (concurrency=%d)", concurrency)
	self.Processors.Start(concurrency, self.StopChannel)
	defer self.Processors.ShutDown()

	if startup != nil {
		go startup()
	}

	<-self.StopChannel

	self.Log.Info("shutting down")

	return nil
}
