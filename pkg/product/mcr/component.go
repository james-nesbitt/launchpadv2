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

// NewMCR constructor for MCR from config.
func NewMCR(id string, c Config) *MCR {
	return &MCR{
		id:     id,
		config: c,
		state:  state{},
	}
}

// MCR product implementation.
type MCR struct {
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
func (p MCR) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, p.id)
}

// Debug product debug.
func (c *MCR) Debug() interface{} {
	return struct {
		Config Config
		State  interface{}
	}{
		Config: c.config,
		State:  c.state.Debug(),
	}
}

// Validate that the project meets the needs of the Product.
func (c *MCR) Validate(ctx context.Context) error {
	errs := []error{}

	if _, err := c.GetManagerHosts(ctx); err != nil {
		errs = append(errs, fmt.Errorf("manager host dependency error: %s", err))
	}
	if _, err := c.GetWorkerHosts(ctx); err != nil {
		errs = append(errs, fmt.Errorf("worker host dependency error: %s", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("MCR validation failed: %s", errors.Join(errs...))
	}
	return nil
}
