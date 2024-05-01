package swarm

import "github.com/Mirantis/launchpad/pkg/implementation/docker"

// NewSwarm constructor.
func NewSwarm(config Config) *Swarm {
	return &Swarm{
		config: config,
	}
}

// Swarm implementation.
type Swarm struct {
	config Config
}

func (s Swarm) ValidateVersion(version docker.Version) error {
	return nil
}
