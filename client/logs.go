package client

import (
	"fmt"
	"io"

	turandotcommon "github.com/tliron/turandot/common"
)

func (self *Client) Logs(appNameSuffix string, containerName string, tail int, follow bool) ([]io.ReadCloser, error) {
	appName := fmt.Sprintf("%s-%s", self.NamePrefix, appNameSuffix)

	if podNames, err := turandotcommon.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		readers := make([]io.ReadCloser, len(podNames))
		for index, podName := range podNames {
			if reader, err := turandotcommon.Log(self.Kubernetes, self.Namespace, podName, containerName, tail, follow); err == nil {
				readers[index] = reader
			} else {
				for i := 0; i < index; i++ {
					readers[i].Close()
				}
				return nil, err
			}
		}
		return readers, nil
	} else {
		return nil, err
	}
}
