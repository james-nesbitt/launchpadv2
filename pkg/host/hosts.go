package host

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// Hosts collection of Host objects for a cluster
type Hosts []Host

// Provides dependency.Dependency for any host type requirements
func (hs Hosts) Provides(ctx context.Context, r dependency.Requirement) error {
	if chr, ok := r.(ClusterHostsRequirement); ok {
		chr.AssignSourceHosts(hs)

		// interupt checking for other requirement handling if validation fails
		if err := r.Validate(ctx); err != nil {
			return err
		}
	}

	return nil // no dependency was handled
}

// --- Dependency on # of hosts of a type of role ---

var (
	ErrRoleHostsDependencyUnMet = fmt.Errorf("%w; cluster could not meet host role dependency requirement", dependency.ErrDependencyNotMet)
)

// ClusterHostsRequirement requirements which use hosts
type ClusterHostsRequirement interface {
	AssignSourceHosts(hs Hosts)
}
