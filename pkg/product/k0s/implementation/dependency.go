package implementation

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ImplementatonType = "k0s-api"
)

type K0sDependencyConfig struct {
	Version string
}

type RequiresK0s interface {
	RequireK0s(context.Context) K0sDependencyConfig
}

type ProvidesK0s interface {
	ProvidesK0s(context.Context) *API
}

func NewRequiresK0s(id string, desc string, config K0sDependencyConfig) *reqK0s {
	return &reqK0s{
		id:     id,
		desc:   desc,
		config: config,
	}
}

func NewK0sDependency(id string, desc string, factory func(context.Context) (*API, error)) *k0sDep {
	return &k0sDep{
		id:      id,
		desc:    desc,
		factory: factory,
	}
}

type reqK0s struct {
	id   string
	desc string

	config K0sDependencyConfig

	dep dependency.Dependency
}

// Id unique identifier for the requirement.
func (r reqK0s) Id() string {
	return r.id
}

// Descibe the requirement.
func (r reqK0s) Describe() string {
	return r.desc
}

// Match the requirement to a dependency.
func (r *reqK0s) Match(d dependency.Dependency) error {
	if _, ok := d.(ProvidesK0s); !ok {
		return fmt.Errorf("%w: k0s did not receive a ProvidesK0s dependency: %+v", dependency.ErrDependencyWrongType, d)
	}

	r.dep = d
	return nil
}

// Matched return the matched dependency, or nil if none has been matched.
func (r *reqK0s) Matched(_ context.Context) dependency.Dependency {
	return r.dep
}

// RequiresK0S indicated that we require k0s (API).
func (r reqK0s) RequiresK0s(_ context.Context) K0sDependencyConfig {
	return r.config
}

// KubernetesDependency Kubernetes providing Dependency.
type K0sDependency interface {
	K0s(ctx context.Context) *API
}

type k0sDep struct {
	id   string
	desc string

	factory func(context.Context) (*API, error)
}

// Id unique identifier for the dependency.
func (d k0sDep) Id() string {
	return d.id
}

// Descibe the dependency.
func (d k0sDep) Describe() string {
	return d.desc
}

// Validate the dependency.
func (d k0sDep) Validate(_ context.Context) error {
	return nil
}

// Met has the dependency been met by the fullfilling backend.
func (d k0sDep) Met(ctx context.Context) error {
	_, err := d.factory(ctx)
	return err
}

// ProvidesK0s fulfill the K0s API dependency.
func (d k0sDep) ProvidesK0s(ctx context.Context) (*API, error) {
	return d.factory(ctx)
}
