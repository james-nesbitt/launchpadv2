package mke4

import (
	"context"
	"fmt"

	kube_rest "k8s.io/client-go/rest"

	"github.com/Mirantis/launchpad/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

func (c *Component) GetKubernetesImplementation(ctx context.Context) (*kubernetes.Kubernetes, error) {
	r := c.k8sr

	if r == nil {
		return nil, fmt.Errorf("%w; no kubernetes requirement created", dependency.ErrRequirementNotCreated)
	}

	d := r.Matched(ctx)
	if d == nil {
		return nil, fmt.Errorf("%w; %s kubernetes requirement not matched", dependency.ErrDependencyNotMatched, r.Id())
	}

	kd, ok := d.(kubernetes.KubernetesDependency) // check that we have a kubernetes dependency
	if !ok {
		// this should never happen, but it is possible
		return nil, fmt.Errorf("%w; %s Dependency is the wrong type for requirement %s", dependency.ErrDependencyNotMatched, d.Id(), r.Id())
	}

	return kd.Kubernetes(ctx)
}

func (c *Component) GetKubernetesRestClient(ctx context.Context) (*kube_rest.RESTClient, error) {
	k, kerr := c.GetKubernetesImplementation(ctx)
	if kerr != nil {
		return nil, kerr
	}

	rc, rcerr := k.RestClient()
	if rcerr != nil {
		return nil, rcerr
	}

	return rc, nil
}
