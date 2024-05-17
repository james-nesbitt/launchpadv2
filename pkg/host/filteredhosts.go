package host

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

var (
	ErrHostsFilterInsufficient = errors.New("not enough hosts in filter")
)

type HostsFilter func(context.Context, Hosts) (Hosts, error)

func HostsDependencyFilterFactory(hs Hosts, filter HostsFilter) func(context.Context) (Hosts, error) {
	return func(ctx context.Context) (Hosts, error) {
		hs, err := filter(ctx, hs)
		if err != nil {
			return hs, fmt.Errorf("Hosts filter produced an error: %s", err.Error())
		}

		return hs, nil
	}
}

func NewHostsFilterRequirement(id string, description string, filter HostsFilter) *hostsFilterRequirement {
	return &hostsFilterRequirement{
		id:     id,
		des:    description,
		filter: filter,
	}
}

type RequiresFilteredHosts interface {
	HostsFilter(context.Context) HostsFilter
}

type hostsFilterRequirement struct {
	id  string
	des string

	filter HostsFilter

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
func (hrr hostsFilterRequirement) HostsFilter(context.Context) HostsFilter {
	return hrr.filter
}
