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
func (c Component) Name() string {
	if c.id == ComponentType {
		return c.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, c.id)
}

// Debug product debug.
func (_ Component) Debug() any { //nolint:staticcheck
	return nil
}

// Validate product debug.
func (_ Component) Validate(_ context.Context) error { //nolint:staticcheck
	return nil
}
