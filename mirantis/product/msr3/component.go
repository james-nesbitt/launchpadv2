package msr3

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "msr3"
)

// NewComponent constructor for MSR3 from config.
func NewComponent(id string, c Config) *Component {
	return &Component{
		id:     id,
		config: c,
		state:  State{},
	}
}

// Component product implementation.
type Component struct {
	id     string
	config Config
	state  State

	k8sr dependency.Requirement
}

// Name for the component.
func (p Component) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return p.id
}

// Debug product debug.
func (_ Component) Debug() interface{} {
	return nil
}

// Validate that the project meets the needs of the Product.
func (_ Component) Validate(context.Context) error {
	return nil
}
