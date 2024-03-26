package msr2

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/phase"
)

// NewMSR2 constructor for MSR2 from config.
func NewMSR2(c Config) MSR2 {
	return MSR2{
		config: c,
		state:  State{},
	}
}

// MSR2 product implementation.
type MSR2 struct {
	config Config
	state  State
}

// Name for the component
func (_ MSR2) Name() string {
	return "MSR2"
}

// Debug product debug.
func (_ MSR2) Debug() interface{} {
	return nil
}

// Provides a list of string labels which this Component provides,
func (_ MSR2) Provides() []string {
	return []string{}
}

// Validate that the cluster meets the needs of the Product
func (_ MSR2) Validate(context.Context, host.Hosts, component.Components) ([]string, error) {
	return []string{}, nil
}

// Actions to run for a particular string phase
func (_ MSR2) Actions(context.Context, string) (phase.Actions, error) {
	return phase.Actions{}, nil
}
