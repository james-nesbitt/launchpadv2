package mke4

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/phase"
)

// NewMKE4 constructor for MKE4 from config.
func NewMKE4(c Config) MKE4 {
	return MKE4{
		config: c,
		state:  State{},
	}
}

// MKE4 product implementation.
type MKE4 struct {
	config Config
	state  State
}

// Name for the component
func (_ MKE4) Name() string {
	return "MKE4"
}

// Debug product debug.
func (_ MKE4) Debug() interface{} {
	return nil
}

// Provides a list of string labels which this Component provides,
func (_ MKE4) Provides() []string {
	return []string{}
}

// Validate that the cluster meets the needs of the Product
func (_ MKE4) Validate(context.Context, host.Hosts, component.Components) ([]string, error) {
	return []string{}, nil
}

// Actions to run for a particular string phase
func (_ MKE4) Actions(context.Context, string) (phase.Actions, error) {
	return phase.Actions{}, nil
}
