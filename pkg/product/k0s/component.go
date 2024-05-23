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

// NewComponent constructor for K0S from config.
func NewComponent(id string, c Config) *Component {
	return &Component{
		id:     id,
		config: c,
		state:  State{},
	}
}

// Component product implementation.
type Component struct {
	id string

	// controller and worker hosts req (hosts onto which we will install)
	hs dependency.Requirement

	// dependencies given out to any component that says that it needs access to our kubernetes/k0s-api
	k8sd dependency.Dependency
	k0sd dependency.Dependency

	config Config
	state  State
}

// Name for the component.
func (p Component) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, p.id)
}

// Debug product debug.
func (_ Component) Debug() interface{} {
	return nil
}

// Validate product debug.
func (_ Component) Validate(_ context.Context) error {
	return nil
}
