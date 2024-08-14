package kubernetes

import (
	"fmt"

	kubeclientcmd "k8s.io/client-go/tools/clientcmd"
)

// Config for the kubernetes implementation
type Config struct {
	KubeCmdApiConfig kubeclientcmd.OverridingClientConfig
}

// ConfigFromKubeConf create a kubernetes implementation configuration from kubeconf bytes
func ConfigFromKubeConf(cb []byte) (Config, error) {
	c := Config{}

	if kckrc, err := kubeclientcmd.NewClientConfigFromBytes(cb); err == nil {
		c.KubeCmdApiConfig = kckrc
	} else {
		return c, fmt.Errorf("could not interpret kubeconfig")
	}

	return c, nil
}
