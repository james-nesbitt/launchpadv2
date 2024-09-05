package kubernetes

import (
	kubeversion "k8s.io/apimachinery/pkg/version"
)

const (
	ImplementationType = "kubernetes"
)

type Version struct {
	Version string
	APIs    []string
}

// NewKubernetes constructor.
func NewKubernetes(config Config) *Kubernetes {
	return &Kubernetes{
		config: config,
	}
}

// Kubernetes implementation.
type Kubernetes struct {
	config Config
}

func (k Kubernetes) Provider() string {
	return k.config.Provider
}

// IsValidaKubernetesVersion is the k8s version an acceptable value for this code base.
func IsValidKubernetesVersion(_ Version) error {
	// @TODO we need to figure out how to perform validation as a feature
	return nil
}

// Determine the ServerVersion (good to test connection).
func (k Kubernetes) ServerVersion() (*kubeversion.Info, error) {
	krc, krcerr := k.Client()
	if krcerr != nil {
		return nil, krcerr
	}

	krcd := krc.DiscoveryClient
	ksv, ksverr := krcd.ServerVersion()
	if ksverr != nil {
		return nil, ksverr
	}

	return ksv, nil
}
