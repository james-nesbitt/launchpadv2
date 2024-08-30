package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

/**
 * Use this for testing cases where you need to test dependency collection,
 * but in most cases the mock.Component is what you want to use
 */

// ProvidesDependencies return a mock function which behaves like a dependency.ProvidesDependencies.
func ProvidesDependencies(dep dependency.Dependency, err error) fillDep {
	return fillDep{
		d:   dep,
		err: err,
	}
}

type fillDep struct {
	d   dependency.Dependency
	err error
}

// ProvidesDependencies provide a dependency for a matching requirement.
func (fd fillDep) ProvidesDependencies(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	return fd.d, fd.err
}
