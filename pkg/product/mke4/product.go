package mke4

import (
	"context"
)

const (
	ComponentType = "mke4"
)

// NewMKE4 constructor for MKE4 from config.
func NewMKE4(id string, c Config) *MKE4 {
	return &MKE4{
		id:     id,
		config: c,
		state:  State{},
	}
}

// MKE4 product implementation.
type MKE4 struct {
	id     string
	config Config
	state  State
}

// Name for the component
func (p MKE4) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return p.id
}

// Debug product debug.
func (_ MKE4) Debug() interface{} {
	return nil
}

// // Validate that the cluster meets the needs of the Product
func (_ MKE4) Validate(context.Context) error {
	return nil
}
