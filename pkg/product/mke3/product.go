package mke3

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/product/mcr"
)

const (
	ComponentType = "mke3"
)

func init() {
	// Add MKE roles to the MCR role list as we need MCR on all nodes
	mcr.MCRHostRoles = append(mcr.MCRHostRoles, []string{"manager", "worker"}...)
}

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

	dhr dependency.Requirement
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
