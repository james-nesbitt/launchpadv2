package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/phase"
)

func NewComponent(name string, action phase.Actions) comp {
	return comp{
		name:    name,
		actions: action,
	}
}

type comp struct {
	name    string
	actions phase.Actions
}

func (c comp) Name() string {
	return c.name
}

// Actions to run for a particular string phase
func (c comp) Actions(context.Context, string) (phase.Actions, error) {
	return c.actions, nil
}

func NewComponentWDependencies(name string, action phase.Actions, requirements dependency.Requirements, dependency dependency.Dependency, providesError error) compWDependencies {
	return compWDependencies{
		comp: comp{
			name:    name,
			actions: action,
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
