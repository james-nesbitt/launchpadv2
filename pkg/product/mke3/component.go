package mke3

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "mke3"
)

// NewMKE3 constructor for MKE3 from config.
func NewMKE3(id string, c Config) *MKE3 {
	return &MKE3{
		id:     id,
		config: c,
		state:  State{},
	}
}

// MKE3 product implementation.
type MKE3 struct {
	id     string
	config Config
	state  State

	dhr dependency.Requirement // Docker hosts dependency

	k8sd dependency.Dependency
	mked dependency.Dependency
}

// Name for the component
func (p MKE3) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return p.id
}

// Debug product debug.
func (_ MKE3) Debug() interface{} {
	return nil
}

// Validate that the cluster meets the needs of the Product
func (_ MKE3) Validate(context.Context) error {
	return nil
}
