package msr2

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "msr2"
)

// NewMSR2 constructor for MSR2 from config.
func NewMSR2(id string, c Config) *MSR2 {
	return &MSR2{
		id:     id,
		config: c,
		state:  State{},
	}
}

// MSR2 product implementation.
type MSR2 struct {
	id     string
	config Config
	state  State

	dhr   dependency.Requirement // docker hosts requirement
	mke3r dependency.Requirement // mke3 requirement

	msr2 dependency.Dependency
}

// Name for the component
func (p MSR2) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return p.id
}

// Debug product debug.
func (_ MSR2) Debug() interface{} {
	return nil
}

// Validate that the cluster meets the needs of the Product
func (_ MSR2) Validate(context.Context) error {
	return nil
}
