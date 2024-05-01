package registry

import "context"

// ProvidesImplementation indicates that it can create a Registry ImplementationRegistry.
type ProvidesImplementation interface {
	ImplementRegistry(context.Context) (Registry, error)
}

// NewRegistry Implementation.
func NewRegistry() Registry {
	return Registry{}
}

// Registry host level Registry implementation.
type Registry struct{}
