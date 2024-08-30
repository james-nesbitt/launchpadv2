package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// helpers

// ReqFactoryStatic return a requirements factory that just returns a static set of requirements.
func ReqFactoryStatic(rs dependency.Requirements) func(context.Context) dependency.Requirements {
	return func(ctx context.Context) dependency.Requirements {
		return rs
	}
}

// DepFactoryStatic return a dependency factory that just returns a static set of dependencies.
func DepFactoryStatic(d dependency.Dependency, err error) func(context.Context, dependency.Requirement) (dependency.Dependency, error) {
	return func(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
		return d, err
	}
}

// DepFactoryIdMatch return a dependency factory that returns a dependendy set if the passed requirement id matches the passed id.
func DepFactoryIdMatch(rid string, d dependency.Dependency) func(context.Context, dependency.Requirement) (dependency.Dependency, error) {
	return func(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
		if rid == r.Id() {
			return d, nil
		}
		return nil, nil
	}
}
