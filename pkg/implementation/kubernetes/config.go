package kubernetes

import (
	"fmt"

	"gopkg.in/yaml.v3"
	kubeclientcmd "k8s.io/client-go/tools/clientcmd"
)

// Config for the kubernetes implementation.
type Config struct {
	KubeCmdApiConfig kubeclientcmd.OverridingClientConfig
}

// ConfigFromKubeConf create a kubernetes implementation configuration from kubeconf bytes.
func ConfigFromKubeConf(cb []byte) (Config, error) {
	var c Config

	if kckrc, err := kubeclientcmd.NewClientConfigFromBytes(cb); err == nil {
		c.KubeCmdApiConfig = kckrc
	} else {
		return c, fmt.Errorf("could not interpret kubeconfig")
	}

	return c, nil
}

// ConfigFromEnv create a kubernetes implementation config object from the environment, using client defaults, with optional overrides
//
//	@NOTE: this should behave similar to how kubectl would, looking for config defaults from ENVVars and default files
func ConfigFromEnv(o *kubeclientcmd.ConfigOverrides) (Config, error) {
	var c Config

	apic, apicerr := kubeclientcmd.NewDefaultPathOptions().GetStartingConfig()
	if apicerr != nil {
		return c, apicerr
	}

	c.KubeCmdApiConfig = kubeclientcmd.NewDefaultClientConfig(*apic, o)

	return c, nil
}

// KubeConfig produce yaml bytes for a kubeconfig for this implementation.
func (c Config) KubeConfig() []byte {
	rc, _ := c.KubeCmdApiConfig.MergedRawConfig()

	mb, mberr := yaml.Marshal(rc.DeepCopy())
	if mberr != nil {
		return []byte{}
	}

	return mb
}
