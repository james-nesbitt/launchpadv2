package k0s

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "k0s"

	// HostRoles roles that MCR considers targets for installation on.
	ControllerHostRole = "controller"
	WorkerHostRole     = "worker"
)

// NewK0S constructor for K0S from config.
func NewK0S(id string, c Config) *K0S {
	return &K0S{
		id:     id,
		config: c,
		state:  State{},
	}
}

// K0S product implementation.
type K0S struct {
	id string

	// controller and worker host roles req (hosts onto which we will install)
	chr dependency.Requirement
	whr dependency.Requirement

	// dependencies given out to any component that says that it needs access to our kubernetes/k0s-api
	k8sd dependency.Dependency
	k0sd dependency.Dependency

	config Config
	state  State
}

// Name for the component.
func (p K0S) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, p.id)
}

// Debug product debug.
func (_ K0S) Debug() interface{} {
	return nil
}

// Validate product debug.
func (_ K0S) Validate(_ context.Context) error {
	return nil
}
