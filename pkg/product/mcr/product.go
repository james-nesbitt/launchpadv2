package mcr

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
	"github.com/Mirantis/launchpad/pkg/implementation/docker/swarm"
	"github.com/Mirantis/launchpad/pkg/phase"
)

// NewMCR constructor for MCR from config.
func NewMCR(c Config) MCR {
	return MCR{
		config: c,
		state:  State{},
	}
}

// MCR product implementation.
type MCR struct {
	config Config
	state  State
}

// Name for the component
func (_ MCR) Name() string {
	return "MCR"
}

// Debug product debug.
func (_ MCR) Debug() interface{} {
	return nil
}

// Provides a list of string labels which this Component provides,
func (_ MCR) Provides() []string {
	return []string{}
}

// Validate that the cluster meets the needs of the Product
func (_ MCR) Validate(context.Context, host.Hosts, component.Components) ([]string, error) {
	return []string{}, nil
}

// Actions to run for a particular string phase
func (_ MCR) Actions(context.Context, string) (phase.Actions, error) {
	return phase.Actions{}, nil
}

// ImplementDocker
func (_ MCR) ImplementDocker(context.Context) (docker.Docker, error) {
	return docker.NewDocker(), nil
}

// ImplementDockerSwarm
func (_ MCR) ImplementDockerSwarm(context.Context) (swarm.Swarm, error) {
	return swarm.NewSwarm(), nil
}
