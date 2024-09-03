package msr4

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

func (c *Component) RequiresDependencies(ctx context.Context) dependency.Requirements {
	if c.k8sr == nil {
		c.k8sr = kubernetes.NewKubernetesRequirement(
			fmt.Sprintf("%s:%s", c.Name(), kubernetes.ImplementationType),
			"MSR4 needs a kubernetes project to install on top of",
			kubernetes.Version{},
		)
	}

	return dependency.Requirements{
		c.k8sr,
	}
}
