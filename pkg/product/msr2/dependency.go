package msr2

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
	mke3implementation "github.com/Mirantis/launchpad/pkg/product/mke3/implementation"
)

// Requires declare that we need a HostsRoles dependency
func (c *MSR2) Requires(_ context.Context) dependency.Requirements {
	if c.dhr == nil {
		c.dhr = dockerhost.NewDockerHostsRequirement(
			fmt.Sprintf("%s:%s", c.Name(), dockerhost.ImplementationType),
			fmt.Sprintf("%s: Needs docker hosts to install on", ComponentType),
			docker.Version{}, // <--- Put in here some limitaitons on what docker version are acceptable
		)
	}
	if c.mke3r == nil {
		c.mke3r = mke3implementation.NewMKE3Requirement(
			fmt.Sprintf("%s:%s", c.Name(), mke3implementation.ImplementationType),
			fmt.Sprintf("%s: relies on MKE for enzi", ComponentType),
			mke3implementation.MKE3DependencyConfig{}, // <--- Put in here some limitations on what version of MKE are acceptable
		)
	}

	return dependency.Requirements{
		c.dhr,
		c.mke3r,
	}
}
