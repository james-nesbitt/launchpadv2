package mke4

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

func (c *MKE4) RequiresDependencies(ctx context.Context) dependency.Requirements {
	if c.k8sr == nil {
		c.k8sr = kubernetes.NewKubernetesRequirement(
			fmt.Sprintf("%s:%s", c.Name(), kubernetes.ImplementationType),
			"MKE4 needs a kubernetes cluster to install on top of",
			kubernetes.Version{},
		)
	}

	return dependency.Requirements{
		c.k8sr,
	}
}
