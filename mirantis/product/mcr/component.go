/*
Package mcr launchpad component plugin for MCR swarm clusters.
*/
package mcr

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "mcr"
)

// NewComponent constructor for MCR from config.
func NewComponent(id string, c Config) *Component {
	return &Component{
		id:     id,
		config: c,
		state:  state{},
	}
}

// Component product implementation.
type Component struct {
	id string

	// host roles req
	hr dependency.Requirement // mcr hosts requirement

	// dependencies given out
	dhd dependency.Dependency // docker (host) dependency
	dsd dependency.Dependency // docker-swarm dependency

	config Config
	state  state
}

// Name for the component.
func (c Component) Name() string {
	if c.id == ComponentType {
		return c.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, c.id)
}

// Debug product debug.
func (c *Component) Debug() any {
	return struct {
		Name   string
		Config Config
		State  any
	}{
		Name:   c.Name(),
		Config: c.config,
		State:  c.state.Debug(),
	}
}

// Validate that the project meets the needs of the Product.
func (c *Component) Validate(ctx context.Context) error {
	errs := []error{}

	if _, err := c.GetManagerHosts(ctx); err != nil {
		errs = append(errs, fmt.Errorf("manager host dependency error: %s", err))
	}
	if _, err := c.GetWorkerHosts(ctx); err != nil {
		errs = append(errs, fmt.Errorf("worker host dependency error: %s", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("mcr validation failed: %s", errors.Join(errs...))
	}
	return nil
}
