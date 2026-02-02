package msr3

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

func (c Component) RequiresDependencies(_ context.Context) (rs dependency.Requirements) {
	if c.k8sr == nil {
		c.k8sr = kubernetes.NewKubernetesRequirement(
			fmt.Sprintf("%s:%s", c.Name(), kubernetes.ImplementationType),
			fmt.Sprintf("MSR3 '%s' needs kubernetes implementation for installation", c.id),
			kubernetes.Version{},
		)
	}
	rs = append(rs, c.k8sr)

	return
}
