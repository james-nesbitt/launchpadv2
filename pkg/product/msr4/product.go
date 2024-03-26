package msr4

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/phase"
)

// NewMSR4 constructor for MSR4 from config.
func NewMSR4(c Config) MSR4 {
	return MSR4{
		config: c,
		state:  State{},
	}
}

// MSR4 product implementation.
type MSR4 struct {
	config Config
	state  State
}

// Name for the component
func (_ MSR4) Name() string {
	return "MSR4"
}

// Debug product debug.
func (_ MSR4) Debug() interface{} {
	return nil
}

// Provides a list of string labels which this Component provides,
func (_ MSR4) Provides() []string {
	return []string{}
}

// Validate that the cluster meets the needs of the Product
func (_ MSR4) Validate(context.Context, host.Hosts, component.Components) ([]string, error) {
	return []string{}, nil
}

// Actions to run for a particular string phase
func (_ MSR4) Actions(context.Context, string) (phase.Actions, error) {
	return phase.Actions{}, nil
}
