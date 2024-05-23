package msr2

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "msr2"
)

// NewComponent constructor for MSR2 from config.
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

	dhr   dependency.Requirement // docker hosts requirement
	mke3r dependency.Requirement // mke3 requirement

	msr2 dependency.Dependency
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
