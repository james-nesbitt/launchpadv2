package k0s

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

// ValidateK8sDependencyConfig validate a Kubernetes client request configuration.
func (c Component) ValidateK8sDependencyConfig(kc kubernetes.Version) error {
	return nil
}

func (c *Component) kubernetesImplementation(ctx context.Context) (*kubernetes.Kubernetes, error) {
	lh := c.GetLeaderHost(ctx)
	if lh == nil {
		return nil, fmt.Errorf("No leader found")
	}

	lkh := HostGetK0s(lh)
	if lkh == nil {
		return nil, fmt.Errorf("leader host has no k0s functionality associated")
	}

	kcs, kcserr := lkh.K0sKubeconfigAdmin(ctx)
	if kcserr != nil {
		return nil, fmt.Errorf("%s: Error retrieving kubeconfig for admin on leader; %w", lh.Id(), kcserr)
	}

	kc, kcerr := kubernetes.ConfigFromKubeConf([]byte(kcs))
	if kcerr != nil {
		return nil, fmt.Errorf("%s: Error interpreting kubeconfig; %w", lh.Id(), kcserr)
	}

	k := kubernetes.NewKubernetes(kc)
	return k, nil
}
