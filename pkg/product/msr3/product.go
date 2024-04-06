package msr3

import (
	"context"
)

const (
	ComponentType = "msr3"
)

// NewMSR3 constructor for MSR3 from config.
func NewMSR3(id string, c Config) *MSR3 {
	return &MSR3{
		id:     id,
		config: c,
		state:  State{},
	}
}

// MSR3 product implementation.
type MSR3 struct {
	id     string
	config Config
	state  State
}

// Name for the component
func (p MSR3) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return p.id
}

// Debug product debug.
func (_ MSR3) Debug() interface{} {
	return nil
}

// Validate that the cluster meets the needs of the Product
func (_ MSR3) Validate(context.Context) error {
	return nil
}
