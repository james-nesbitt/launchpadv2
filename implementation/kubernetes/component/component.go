package kubernetes

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "kubernetes"
)

// NewComponent constructor for K0S from config.
func NewComponent(id string, c kubernetes.Config) *Component {
	return &Component{
		id:     id,
		config: c,
		state:  State{}, // still not doing anything with state :(
	}
}

// Component product implementation.
type Component struct {
	id string

	// dependencies given out to any component that says that it needs access to our kubernetes
	k8sd dependency.Dependency

	config kubernetes.Config
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
