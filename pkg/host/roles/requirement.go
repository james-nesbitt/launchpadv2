package roles

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

/**
 * Requirement for Hosts that match roles
 *
 * Use this to require hosts that match roles from a
 * dependency provider that can match roles.
 */

var (
	RequiresHostsRolesType = "host-roles"
)

func NewHostsRolesRequirement(id string, description string, filter HostsRolesFilter) *hostsRolesRequirement {
	return &hostsRolesRequirement{
		id:     id,
		des:    description,
		filter: filter,
	}
}

type RequiresHostsRoles interface {
	RequireHostsRoles(context.Context) HostsRolesFilter
}

type hostsRolesRequirement struct {
	id  string
	des string

	filter HostsRolesFilter

	dep dependency.Dependency
}

// Requirer identify the requirer with a string, which may be used for correlation in the future.
func (hrr hostsRolesRequirement) Id() string {
	return hrr.id
}

// Describe the Requirements for labelling/auditing.
func (hrr hostsRolesRequirement) Describe() string {
	return hrr.des
}

// Match with a Dependency.
func (hrr *hostsRolesRequirement) Match(dep dependency.Dependency) error {
	hrr.dep = dep
	return nil
}

// Matched has been Matched. If not matched return nil.
func (hrr hostsRolesRequirement) Matched(_ context.Context) dependency.Dependency {
	return hrr.dep
}

// RequiresHostsRoles.
func (hrr hostsRolesRequirement) RequireHostsRoles(context.Context) HostsRolesFilter {
	return hrr.filter
}
