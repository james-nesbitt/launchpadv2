package msr3

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

func (p Component) RequiresDependencies(_ context.Context) (rs dependency.Requirements) {
	if p.k8sr == nil {
		p.k8sr = kubernetes.NewKubernetesRequirement(
			fmt.Sprintf("%s:%s", p.Name(), kubernetes.ImplementationType),
			fmt.Sprintf("MSR3 '%s' needs kubernetes implementation for installation", p.id),
			kubernetes.Version{},
		)
	}
	rs = append(rs, p.k8sr)

	return
}
