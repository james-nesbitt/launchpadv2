package mcr

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "mcr"
)

var (
	// MCRManagerHostRoles the Host roles accepted for managers.
	MCRManagerHostRoles = []string{"manager"}
	// MCRWorkerHostRoles the Host roles accepted for workers.
	MCRWorkerHostRoles = []string{"worker"}
)

// NewMCR constructor for MCR from config.
func NewMCR(id string, c Config) *MCR {
	return &MCR{
		id:     id,
		config: c,
		state:  state{hosts: map[string]*hostState{}},
	}
}

// MCR product implementation.
type MCR struct {
	id string

	// host roles req
	mhr dependency.Requirement // manager hosts requirement
	whr dependency.Requirement // worker hosts requirement

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
func (_ MCR) Debug() interface{} {
	return nil
}

// Validate that the cluster meets the needs of the Product.
func (_ MCR) Validate(context.Context) error {
	return nil
}
