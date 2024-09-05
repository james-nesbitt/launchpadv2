package kubernetes

import (
	"fmt"

	helm "github.com/mittwald/go-helm-client"
	kube_discovery "k8s.io/client-go/discovery"
	kube_dynamic "k8s.io/client-go/dynamic"
	kube_kubernetes "k8s.io/client-go/kubernetes"
	kube_rest "k8s.io/client-go/rest"
)

// -- Kube client generation ---

// DiscoveryClient kubernetes discovery client
func (k Kubernetes) DiscoveryClient() (*kube_discovery.DiscoveryClient, error) {
	krc, krcerr := k.Client()
	if krcerr != nil {
		return nil, krcerr
	}

	return krc.DiscoveryClient, nil
}

// Client go-client kubernetes.Interface
func (k Kubernetes) Client() (*kube_kubernetes.Clientset, error) {
	if k.config.KubeCmdApiConfig == nil {
		return nil, fmt.Errorf("kubernetes implementation has no valid config options")
	}

	rc, rcerr := k.config.KubeCmdApiConfig.ClientConfig()
	if rcerr != nil {
		return nil, rcerr
	}
	cl, clerr := kube_kubernetes.NewForConfig(rc)
	if clerr != nil {
		return nil, clerr
	}

	return cl, nil
}

// DynamicClient create and return a kubernetes dynamic client from the config
func (k Kubernetes) DynamicClient() (*kube_dynamic.DynamicClient, error) {
	rcc, rccerr := k.config.KubeCmdApiConfig.ClientConfig()
	if rccerr != nil {
		return nil, rccerr
	}

	dc, dcerr := kube_dynamic.NewForConfig(rcc)
	if dcerr != nil {
		return nil, dcerr
	}

	return dc, nil
}

// HelmClient for the kubernetes implementation.
func (k Kubernetes) HelmClient(opts helm.Options) (helm.Client, error) {
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

// RestClient create the kubernetes rest client from the config
func (k Kubernetes) RestClient() (*kube_rest.RESTClient, error) {
	rcc, rccerr := k.config.KubeCmdApiConfig.ClientConfig()
	if rccerr != nil {
		return nil, rccerr
	}

	rc, rcerr := kube_rest.RESTClientFor(rcc)
	if rcerr != nil {
		return nil, rcerr
	}

	return rc, nil
}
