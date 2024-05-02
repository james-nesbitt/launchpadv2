package implementation

import (
	"context"
)

// ProvidesK0s can provide the k0s API implementation.
type ProvidesK0s interface {
	ProvidesK0s(context.Context) *API
}

// NewK0sDependency build a new K0S providing Dependency.
func NewK0sDependency(id string, desc string, factory func(context.Context) (*API, error)) *k0sDep {
	return &k0sDep{
		id:      id,
		desc:    desc,
		factory: factory,
	}
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
