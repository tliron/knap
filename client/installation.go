package client

import (
	"fmt"
	"time"

	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	"github.com/tliron/turandot/version"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	waitpkg "k8s.io/apimachinery/pkg/util/wait"
)

func (self *Client) Install(registry string, wait bool) error {
	var err error

	if _, err = self.createCustomResourceDefinition(); err != nil {
		return err
	}

	if _, err = self.createNamespace(); err != nil {
		return err
	}

	var serviceAccount *core.ServiceAccount
	if serviceAccount, err = self.createServiceAccount(); err != nil {
		return err
	}

	if self.Cluster {
		if _, err = self.createClusterRoleBinding(serviceAccount); err != nil {
			return err
		}
	} else {
		var role *rbac.Role
		if role, err = self.createRole(); err != nil {
			return err
		}
		if _, err = self.createRoleBinding(serviceAccount, role); err != nil {
			return err
		}
	}

	var operatorDeployment *apps.Deployment
	if operatorDeployment, err = self.createOperatorDeployment(registry, serviceAccount, 1); err != nil {
		return err
	}

	if wait {
		if _, err := self.waitForDeployment(operatorDeployment.Name); err != nil {
			return err
		}
	}

	return nil
}

func (self *Client) Uninstall(wait bool) {
	var gracePeriodSeconds int64 = 0
	deleteOptions := meta.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
	}

	// Operator deployment
	if err := self.Kubernetes.AppsV1().Deployments(self.Namespace).Delete(self.Context, fmt.Sprintf("%s-operator", self.NamePrefix), deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	if self.Cluster {
		// Cluster role binding
		if err := self.Kubernetes.RbacV1().ClusterRoleBindings().Delete(self.Context, self.NamePrefix, deleteOptions); err != nil {
			self.Log.Warningf("%s", err)
		}
	} else {
		// Role binding
		if err := self.Kubernetes.RbacV1().RoleBindings(self.Namespace).Delete(self.Context, self.NamePrefix, deleteOptions); err != nil {
			self.Log.Warningf("%s", err)
		}

		// Role
		if err := self.Kubernetes.RbacV1().Roles(self.Namespace).Delete(self.Context, self.NamePrefix, deleteOptions); err != nil {
			self.Log.Warningf("%s", err)
		}
	}

	// Service account
	if err := self.Kubernetes.CoreV1().ServiceAccounts(self.Namespace).Delete(self.Context, self.NamePrefix, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Custom resource definition
	if err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Delete(self.Context, resources.NetworkCustomResourceDefinition.Name, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	if wait {
		self.waitForDeletion("service", func() bool {
			_, err := self.Kubernetes.CoreV1().Services(self.Namespace).Get(self.Context, fmt.Sprintf("%s-inventory", self.NamePrefix), meta.GetOptions{})
			return err == nil
		})
		self.waitForDeletion("operator deployment", func() bool {
			_, err := self.Kubernetes.AppsV1().Deployments(self.Namespace).Get(self.Context, fmt.Sprintf("%s-operator", self.NamePrefix), meta.GetOptions{})
			return err == nil
		})
		if self.Cluster {
			self.waitForDeletion("cluster role binding", func() bool {
				_, err := self.Kubernetes.RbacV1().ClusterRoleBindings().Get(self.Context, self.NamePrefix, meta.GetOptions{})
				return err == nil
			})
		} else {
			self.waitForDeletion("role binding", func() bool {
				_, err := self.Kubernetes.RbacV1().RoleBindings(self.Namespace).Get(self.Context, self.NamePrefix, meta.GetOptions{})
				return err == nil
			})
			self.waitForDeletion("role", func() bool {
				_, err := self.Kubernetes.RbacV1().Roles(self.Namespace).Get(self.Context, self.NamePrefix, meta.GetOptions{})
				return err == nil
			})
		}
		self.waitForDeletion("service account", func() bool {
			_, err := self.Kubernetes.CoreV1().ServiceAccounts(self.Namespace).Get(self.Context, self.NamePrefix, meta.GetOptions{})
			return err == nil
		})
		self.waitForDeletion("custom resource definition", func() bool {
			_, err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.Context, resources.NetworkCustomResourceDefinition.Name, meta.GetOptions{})
			return err == nil
		})
	}
}

func (self *Client) createCustomResourceDefinition() (*apiextensions.CustomResourceDefinition, error) {
	customResourceDefinition := &resources.NetworkCustomResourceDefinition

	if customResourceDefinition, err := self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Create(self.Context, customResourceDefinition, meta.CreateOptions{}); err == nil {
		return customResourceDefinition, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.APIExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.Context, resources.NetworkCustomResourceDefinition.Name, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createNamespace() (*core.Namespace, error) {
	namespace := &core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: self.Namespace,
		},
	}

	if namespace, err := self.Kubernetes.CoreV1().Namespaces().Create(self.Context, namespace, meta.CreateOptions{}); err == nil {
		return namespace, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.CoreV1().Namespaces().Get(self.Context, self.Namespace, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createServiceAccount() (*core.ServiceAccount, error) {
	serviceAccount := &core.ServiceAccount{
		ObjectMeta: meta.ObjectMeta{
			Name: self.NamePrefix,
		},
	}

	if serviceAccount, err := self.Kubernetes.CoreV1().ServiceAccounts(self.Namespace).Create(self.Context, serviceAccount, meta.CreateOptions{}); err == nil {
		return serviceAccount, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.CoreV1().ServiceAccounts(self.Namespace).Get(self.Context, self.NamePrefix, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createRole() (*rbac.Role, error) {
	role := &rbac.Role{
		ObjectMeta: meta.ObjectMeta{
			Name: self.NamePrefix,
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{rbac.APIGroupAll},
				Resources: []string{rbac.ResourceAll},
				Verbs:     []string{rbac.VerbAll},
			},
		},
	}

	if role, err := self.Kubernetes.RbacV1().Roles(self.Namespace).Create(self.Context, role, meta.CreateOptions{}); err == nil {
		return role, err
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.RbacV1().Roles(self.Namespace).Get(self.Context, self.NamePrefix, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createRoleBinding(serviceAccount *core.ServiceAccount, role *rbac.Role) (*rbac.RoleBinding, error) {
	roleBinding := &rbac.RoleBinding{
		ObjectMeta: meta.ObjectMeta{
			Name: self.NamePrefix,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind, // serviceAccount.Kind is empty
				Name:      serviceAccount.Name,
				Namespace: self.Namespace, // required
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName, // role.GroupVersionKind().Group is empty
			Kind:     "Role",         // role.Kind is empty
			Name:     role.Name,
		},
	}

	if roleBinding, err := self.Kubernetes.RbacV1().RoleBindings(self.Namespace).Create(self.Context, roleBinding, meta.CreateOptions{}); err == nil {
		return roleBinding, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.RbacV1().RoleBindings(self.Namespace).Get(self.Context, self.NamePrefix, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createClusterRoleBinding(serviceAccount *core.ServiceAccount) (*rbac.ClusterRoleBinding, error) {
	clusterRoleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: meta.ObjectMeta{
			Name: self.NamePrefix,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind, // serviceAccount.Kind is empty
				Name:      serviceAccount.Name,
				Namespace: self.Namespace, // required
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}

	if clusterRoleBinding, err := self.Kubernetes.RbacV1().ClusterRoleBindings().Create(self.Context, clusterRoleBinding, meta.CreateOptions{}); err == nil {
		return clusterRoleBinding, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.RbacV1().ClusterRoleBindings().Get(self.Context, self.NamePrefix, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createOperatorDeployment(registry string, serviceAccount *core.ServiceAccount, replicas int32) (*apps.Deployment, error) {
	appName := fmt.Sprintf("%s-operator", self.NamePrefix)
	instanceName := fmt.Sprintf("%s-%s", appName, self.Namespace)

	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name: appName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   instanceName,
				"app.kubernetes.io/version":    version.GitVersion,
				"app.kubernetes.io/component":  "operator",
				"app.kubernetes.io/part-of":    self.PartOf,
				"app.kubernetes.io/managed-by": self.ManagedBy,
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":      appName,
					"app.kubernetes.io/instance":  instanceName,
					"app.kubernetes.io/version":   version.GitVersion,
					"app.kubernetes.io/component": "operator",
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":       appName,
						"app.kubernetes.io/instance":   instanceName,
						"app.kubernetes.io/version":    version.GitVersion,
						"app.kubernetes.io/component":  "operator",
						"app.kubernetes.io/part-of":    self.PartOf,
						"app.kubernetes.io/managed-by": self.ManagedBy,
					},
				},
				Spec: core.PodSpec{
					ServiceAccountName: serviceAccount.Name,
					Containers: []core.Container{
						{
							Name:            "operator",
							Image:           fmt.Sprintf("%s/%s", registry, self.OperatorImageName),
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{
								{
									Name:  "KNAP_OPERATOR_concurrency",
									Value: "3",
								},
								{
									Name:  "KNAP_OPERATOR_verbose",
									Value: "1",
								},
							},
							LivenessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{
										Port: intstr.FromInt(8086),
										Path: "/live",
									},
								},
							},
							ReadinessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{
										Port: intstr.FromInt(8086),
										Path: "/ready",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return self.createDeployment(deployment, appName)
}

func (self *Client) createDeployment(deployment *apps.Deployment, appName string) (*apps.Deployment, error) {
	if deployment, err := self.Kubernetes.AppsV1().Deployments(self.Namespace).Create(self.Context, deployment, meta.CreateOptions{}); err == nil {
		return deployment, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.AppsV1().Deployments(self.Namespace).Get(self.Context, appName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) waitForDeletion(name string, condition func() bool) {
	err := waitpkg.PollImmediate(time.Second, timeout, func() (bool, error) {
		self.Log.Infof("waiting for %s to delete", name)
		return !condition(), nil
	})
	if err != nil {
		self.Log.Warningf("error while waiting for %s to delete: %s", name, err.Error())
	}
}
