package k0s

import (
	"context"
	"fmt"
)

const (
	ComponentType = "k0s"
)

// NewK0S constructor for K0S from config.
func NewK0S(id string, c Config) *K0S {
	return &K0S{
		id:     id,
		config: c,
		state:  State{},
	}
}

// K0S product implementation.
type K0S struct {
	id     string
	config Config
	state  State
}

// Name for the component
func (p K0S) Name() string {
	if p.id == ComponentType {
		return p.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, p.id)
}

// Debug product debug.
func (_ K0S) Debug() interface{} {
	return nil
}

// Validate product debug.
func (_ K0S) Validate(_ context.Context) error {
	return nil
}
