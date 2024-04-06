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
	// MCRHostRoles roles that MCR considers targets for installation on
	// Other products like MKE add roles to this list.
	MCRHostRoles = []string{"mcr"}
)

// NewMCR constructor for MCR from config.
func NewMCR(id string, c Config) *MCR {
	return &MCR{
		id:     id,
		config: c,
		state:  State{},
	}
}

// MCR product implementation.
type MCR struct {
	id string

	// host roles req
	rhr dependency.Requirement

	config Config
	state  State
}

// Name for the component
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

// Validate that the cluster meets the needs of the Product
func (_ MCR) Validate(context.Context) error {
	return nil
}
