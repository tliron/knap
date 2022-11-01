package client

import (
	"fmt"

	"github.com/tliron/kutil/version"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var true_ = true
var false_ = false

func (self *Client) CreateDeployment(deployment *apps.Deployment, appName string) (*apps.Deployment, error) {
	if deployment, err := self.Kubernetes.AppsV1().Deployments(self.Namespace).Create(self.Context, deployment, meta.CreateOptions{}); err == nil {
		return deployment, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.AppsV1().Deployments(self.Namespace).Get(self.Context, appName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) Labels(appName string, component string, namespace string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       appName,
		"app.kubernetes.io/instance":   fmt.Sprintf("%s-%s", appName, namespace),
		"app.kubernetes.io/version":    version.GitVersion,
		"app.kubernetes.io/component":  component,
		"app.kubernetes.io/part-of":    self.PartOf,
		"app.kubernetes.io/managed-by": self.ManagedBy,
	}
}

func (self *Client) DefaultSecurityContext() *core.SecurityContext {
	var user int64 = 1000
	return &core.SecurityContext{
		AllowPrivilegeEscalation: &false_,
		Capabilities: &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		},
		RunAsNonRoot: &true_,
		RunAsUser:    &user,
		SeccompProfile: &core.SeccompProfile{
			Type: core.SeccompProfileTypeRuntimeDefault,
		},
	}
}
