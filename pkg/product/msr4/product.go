package msr4

import (
	"context"
)

const (
	ComponentType = "msr4"
)

// NewMSR4 constructor for MSR4 from config.
func NewMSR4(id string, c Config) *MSR4 {
	return &MSR4{
		id:     id,
		config: c,
		state:  State{},
	}
}

// MSR4 product implementation.
type MSR4 struct {
	id     string
	config Config
	state  State
}

// Name for the component.
func (p MSR4) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return p.id
}

// Debug product debug.
func (_ MSR4) Debug() interface{} {
	return nil
}

// Validate that the cluster meets the needs of the Product.
func (_ MSR4) Validate(context.Context) error {
	return nil
}
