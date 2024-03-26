package msr3

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/phase"
)

// NewMSR3 constructor for MSR3 from config.
func NewMSR3(c Config) MSR3 {
	return MSR3{
		config: c,
		state:  State{},
	}
}

// MSR3 product implementation.
type MSR3 struct {
	config Config
	state  State
}

// Name for the component
func (_ MSR3) Name() string {
	return "MSR3"
}

// Debug product debug.
func (_ MSR3) Debug() interface{} {
	return nil
}

// Provides a list of string labels which this Component provides,
func (_ MSR3) Provides() []string {
	return []string{}
}

// Validate that the cluster meets the needs of the Product
func (_ MSR3) Validate(context.Context, host.Hosts, component.Components) ([]string, error) {
	return []string{}, nil
}

// Actions to run for a particular string phase
func (_ MSR3) Actions(context.Context, string) (phase.Actions, error) {
	return phase.Actions{}, nil
}
