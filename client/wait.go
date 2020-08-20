package client

import (
	"fmt"
	"time"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	waitpkg "k8s.io/apimachinery/pkg/util/wait"
)

var timeout = 60 * time.Second

func (self *Client) waitForDeployment(appName string) (*apps.Deployment, error) {
	self.Log.Infof("waiting for deployment for %q", appName)

	var deployment *apps.Deployment
	err := waitpkg.PollImmediate(time.Second, timeout, func() (bool, error) {
		var err error
		if deployment, err = self.Kubernetes.AppsV1().Deployments(self.Namespace).Get(self.Context, appName, meta.GetOptions{}); err == nil {
			for _, condition := range deployment.Status.Conditions {
				switch condition.Type {
				case apps.DeploymentAvailable:
					if condition.Status == core.ConditionTrue {
						return true, nil
					}
				case apps.DeploymentReplicaFailure:
					if condition.Status == core.ConditionTrue {
						return false, fmt.Errorf("replica failure: %s", appName)
					}
				}
			}
			return false, nil
		} else {
			return false, err
		}
	})

	if err == nil {
		self.Log.Infof("deployment available for %q", appName)
		//if err := self.waitForPods(appName, deployment); err == nil {
		return deployment, nil
		/*} else {
			return nil, err
		}*/
	} else {
		return nil, err
	}
}
