/*
Package msr2 MSR legacy component for launchpad.
*/
package msr2

import (
	"context"
	"fmt"

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
func (c Component) Name() string {
	if c.id == ComponentType {
		return c.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, c.id)
}

// Debug product debug.
func (c Component) Debug() any {
	return fmt.Sprintf("%s '%s' debug", ComponentType, c.Name())
}

// Validate that the project meets the needs of the Product.
func (_ Component) Validate(context.Context) error { //nolint:staticcheck

	return nil
}
