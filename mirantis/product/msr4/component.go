/*
Package msr4 MSR4 component for launchpad.
*/
package msr4

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "msr4"
)

// NewComponent constructor for MSR4 from config.
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
