package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

func NewComponent(name string, debug interface{}, validate error) comp {
	return comp{
		name:     name,
		debug:    debug,
		validate: validate,
	}
}

type comp struct {
	name     string
	debug    interface{}
	validate error
}

func (c comp) Name() string {
	return c.name
}

func (c comp) Debug() interface{} {
	return struct {
		Name string
	}{
		Name: c.name,
	}
}

func (c comp) Validate(_ context.Context) error {
	return c.validate
}

func NewComponentWDependencies(name string, debug interface{}, validate error, requirements dependency.Requirements, dependency dependency.Dependency, providesError error) compWDependencies {
	return compWDependencies{
		comp: comp{
			name:     name,
			debug:    debug,
			validate: validate,
		},
		reqs:     requirements,
		dep:      dependency,
		provides: providesError,
	}
}

type compWDependencies struct {
	comp

	reqs     dependency.Requirements
	dep      dependency.Dependency
	provides error
}

func (cwd compWDependencies) Requires(context.Context) dependency.Requirements {
	return cwd.reqs
}

func (cwd compWDependencies) Provides(context.Context, dependency.Requirement) (dependency.Dependency, error) {
	return cwd.dep, cwd.provides
}
