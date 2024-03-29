package k0s

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/phase"
)

// NewK0S constructor for K0S from config.
func NewK0S(c Config) K0S {
	return K0S{
		config: c,
		state:  State{},
	}
}

// K0S product implementation.
type K0S struct {
	config Config
	state  State
}

// Name for the component
func (_ K0S) Name() string {
	return "K0S"
}

// Debug product debug.
func (_ K0S) Debug() interface{} {
	return nil
}

// Actions to run for a particular string phase
func (_ K0S) Actions(context.Context, string) (phase.Actions, error) {
	return phase.Actions{}, nil
}
