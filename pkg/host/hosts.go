package host

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// Hosts collection of Host objects for a cluster
type Hosts []Host

// Provides dependency.Dependency for any host type requirements
func (hs Hosts) Provides(ctx context.Context, r dependency.Requirement) dependency.Dependency {
	if chr, ok := r.(ClusterHostsRoleRequirement); !ok {
		d := DependencyFromHosts(hs, chr)
		return d // d.Met() can describe any errors
	}

	return nil // no dependency was handled
}

// --- Dependency on # of hosts of a type of role ---

var (
	ErrRoleHostsDependencyUnMet = fmt.Errorf("%w; cluster could not meet host role dependency requirement", dependency.ErrDependencyNotMet)
)

func DependencyFromHosts(hs Hosts, chr ClusterHostsRoleRequirement) dependency.Dependency {
	rhs := Hosts{}

	for _, hs := range hs {
		if chr.MaximumCount > 0 && len(rhs) > chr.MaximumCount {
			break
		}

		if !hs.HasRole(chr.Role) {
			continue
		}

		rhs = append(rhs, hs)
	}

	var err error
	if chr.MinimumCount > 0 && len(rhs) < chr.MinimumCount {
		err = fmt.Errorf("%w; %s", ErrRoleHostsDependencyUnMet, chr.Describe())
	}

	d := ClusterHostsRoleDependency{
		Description: chr.Describe(),
		Hosts:       hs,
		Error:       err,
	}

	return d
}

// ClusterHostsRoleRequirement a requirement for hosts with matching roles
type ClusterHostsRoleRequirement struct {
	Description string
	// Role of hosts needed
	Role string
	// If 0 then no limit - but describe minimum needed hosts, and maximum consumed hosts
	MinimumCount, MaximumCount int

	// dependency if resolved
	dependency dependency.Dependency
}

func (chr ClusterHostsRoleRequirement) Describe() string {
	return chr.Description
}

func (chr ClusterHostsRoleRequirement) Fullfill(d dependency.Dependency) error {
	chr.dependency = d
	return nil
}

type ClusterHostsRoleDependency struct {
	Description string
	Hosts       Hosts
	Error       error
}

func (chrd ClusterHostsRoleDependency) Describe() string {
	return chrd.Description
}

func (chrd ClusterHostsRoleDependency) Met() error {
	return chrd.Error
}
