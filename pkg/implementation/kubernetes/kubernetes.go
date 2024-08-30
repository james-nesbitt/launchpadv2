package kubernetes

import (
	"context"
	"fmt"

	helm "github.com/mittwald/go-helm-client"
	kubeversion "k8s.io/apimachinery/pkg/version"
	kube "k8s.io/client-go/kubernetes"
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

// IsValidaKubernetesVersion is the k8s version an acceptable value for this code base.
func IsValidKubernetesVersion(_ Version) error {
	return nil
}

// KubeRestClientset for the kubernetes implementation.
func (k Kubernetes) KubeRestClientset(ctx context.Context) (*kube.Clientset, error) {
	if k.config.KubeCmdApiConfig == nil {
		return nil, fmt.Errorf("kubernetes implementation has no valid config options")
	}

	rc, rcerr := k.config.KubeCmdApiConfig.ClientConfig()
	if rcerr != nil {
		return nil, rcerr
	}
	cl, clerr := kube.NewForConfig(rc)
	if clerr != nil {
		return nil, clerr
	}

	return cl, nil
}

// HelmClient for the kubernetes implementation.
func (k Kubernetes) HelmClient(ctx context.Context, opts helm.Options) (helm.Client, error) {
	rc, rcerr := k.config.KubeCmdApiConfig.ClientConfig()
	if rcerr != nil {
		return nil, rcerr
	}

	opt := &helm.RestConfClientOptions{
		Options:    &opts,
		RestConfig: rc,
	}

	hc, err := helm.NewClientFromRestConf(opt)
	if err != nil {
		panic(err)
	}

	return hc, nil
}

// Determine the ServerVersion (good to test connection).
func (k Kubernetes) ServerVersion(ctx context.Context) (*kubeversion.Info, error) {
	krc, krcerr := k.KubeRestClientset(ctx)
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
