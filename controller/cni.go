package controller

import (
	"fmt"
	"strings"

	resources "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	"github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/transcribe"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func (self *Controller) createCniConfig(network *resources.Network) (string, error) {
	appName := fmt.Sprintf("knap-provider-%s", network.Spec.Provider)

	podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Client.Namespace, appName)
	if err != nil {
		return "", fmt.Errorf("cannot find provider for network %s/%s: %s\n%s", network.Namespace, network.Name, network.Spec.Provider, err.Error())
	}

	var stdout strings.Builder
	var stderr strings.Builder

	execOptions := core.PodExecOptions{
		Container: "provider",
		Command:   []string{appName, "provide", network.Name},
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}

	streamOptions := remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	}

	if network.Spec.Hints != nil {
		if hints, err := transcribe.EncodeYAML(network.Spec.Hints, "  ", false); err == nil {
			execOptions.Stdin = true
			streamOptions.Stdin = strings.NewReader(hints)
		} else {
			return "", err
		}
	}

	request := self.REST.Post().Namespace(self.Client.Namespace).Resource("pods").Name(podName).SubResource("exec").VersionedParams(&execOptions, scheme.ParameterCodec)

	if executor, err := remotecommand.NewSPDYExecutor(self.Config, "POST", request.URL()); err == nil {
		if err = executor.Stream(streamOptions); err == nil {
			return stdout.String(), nil
		} else {
			return "", fmt.Errorf("%s\n%s", err.Error(), stderr.String())
		}
	} else {
		return "", err
	}
}
