package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "mock"
)

// NewSimpleComponent simplest of mocking components.
func NewSimpleComponent(name string, debug interface{}, validate error) *comp {
	return &comp{
		name:     name,
		debug:    debug,
		validate: validate,
	}
}

// comp simple mocking component.
type comp struct {
	name     string
	debug    interface{}
	validate error
}

// Name of the Component, which is used for logging and auditing.
func (c comp) Name() string {
	return c.name
}

// Debug the component by returning something that can be serialized.
func (c comp) Debug() interface{} {
	return ComponentDebug{
		Name:  c.name,
		Debug: c.debug,
	}
}

// Validate the component thinks it has valid config.
func (c comp) Validate(_ context.Context) error {
	return c.validate
}

// NewComponentWDependencies mocking component with dependencies.
func NewComponentWDependencies(name string, debug interface{}, validate error, reqFactory func(context.Context) dependency.Requirements, depFactory func(context.Context, dependency.Requirement) (dependency.Dependency, error)) *compWDepedencies {
	return &compWDepedencies{
		comp:       comp{},
		reqFactory: reqFactory,
		depFactory: depFactory,
	}
}

// compWDepedencies simple component that also includes dependency handling using factories.
type compWDepedencies struct {
	comp

	reqFactory func(context.Context) dependency.Requirements
	depFactory func(context.Context, dependency.Requirement) (dependency.Dependency, error)
}

// Requires declare requirements needed for this component.
func (c *compWDepedencies) RequiresDependencies(ctx context.Context) dependency.Requirements {
	if c.reqFactory == nil {
		return nil
	}

	return c.reqFactory(ctx)
}

// ProvidesDependencies dependencies that this component provides.
func (c *compWDepedencies) ProvidesDependencies(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if c.depFactory == nil {
		return nil, nil
	}

	return c.depFactory(ctx, r)
}
