package kubernetes

import (
	"context"

	kube "k8s.io/client-go/kubernetes"

	helm "github.com/mittwald/go-helm-client"
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

func (k Kubernetes) KubeRestClientset(ctx context.Context) (*kube.Clientset, error) {
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
