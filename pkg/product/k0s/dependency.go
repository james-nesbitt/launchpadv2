package k0s

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// --- Dependency requirements ---

// --- Dependency provides ---

func (k K0S) Provides(ctx context.Context, r dependency.Requirement) dependency.Dependency {
	if chr, ok := r.(ClusterHostsRoleRequirement); !ok {
		d := DependencyFromHosts(hs, chr)
		return d // d.Met() can describe any errors
	}

	return nil // no dependency was handled

}

type KubernetesImplementationRequirement struct {

}

//// ImplementKubernetes
//func (_ K0S) ImplementKubernetes(context.Context) (kubernetes.Kubernetes, error) {
//	return kubernetes.NewKubernetes(), nil
//}
//
//// ImplementK0sApi
//func (_ K0S) ImplementK0SApi(context.Context) (implementation.API, error) {
//	return implementation.NewAPI(), nil
//}
//
