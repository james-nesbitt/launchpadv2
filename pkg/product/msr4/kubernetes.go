package msr4

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

func (c Component) GetKubernetesImplementation(ctx context.Context) (*kubernetes.Kubernetes, error) {
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
