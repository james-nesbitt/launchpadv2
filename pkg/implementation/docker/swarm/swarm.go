package swarm

import "context"

// ProvidesImplementation indicates that a kubernetes implementation can be created
type ProvidesImplementation interface {
	ImplementDockerSwarm(context.Context) (Swarm, error)
}

// NewSwarm constructor
func NewSwarm() Swarm {
	return Swarm{}
}

// Swarm implementation
type Swarm struct {
}
