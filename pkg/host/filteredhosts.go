package host

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

func NewHostsFilterRequirement(id string, description string, filter func(context.Context, Hosts) (Hosts, error)) *hostsFilterRequirement {
	return &hostsFilterRequirement{
		id:     id,
		des:    description,
		filter: filter,
	}
}

type RequiresFilteredHosts interface {
	RequireFilteredHosts(context.Context) func(context.Context, Hosts) (Hosts, error)
}

type hostsFilterRequirement struct {
	id  string
	des string

	filter func(context.Context, Hosts) (Hosts, error)

	dep dependency.Dependency
}

// Requirer identify the requirer with a string, which may be used for correlation in the future.
func (hrr hostsFilterRequirement) Id() string {
	return hrr.id
}

// Describe the Requirements for labelling/auditing.
func (hrr hostsFilterRequirement) Describe() string {
	return hrr.des
}

// Match with a Dependency.
func (hrr *hostsFilterRequirement) Match(dep dependency.Dependency) error {
	hrr.dep = dep
	return nil
}

// Matched has been Matched. If not matched return nil.
func (hrr hostsFilterRequirement) Matched(_ context.Context) dependency.Dependency {
	return hrr.dep
}

// RequiresHostsFilter.
func (hrr hostsFilterRequirement) RequireHostsFilter(context.Context) func(context.Context, Hosts) (Hosts, error) {
	return hrr.filter
}
