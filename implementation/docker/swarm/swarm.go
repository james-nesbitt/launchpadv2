package swarm

import (
	"context"
	"fmt"

	dockerimplementation "github.com/Mirantis/launchpad/implementation/docker"
)

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

func (s Swarm) ValidateVersion(version dockerimplementation.Version) error {
	return nil
}

// SwarmAddress convert an IP into a address.
func SwarmAddress(ip string) string {
	return fmt.Sprintf("%s:%d", ip, 2377)
}

// IsDockerSwarmManager return is a docker implementation points to a swarm manager.
func IsDockerSwarmManager(ctx context.Context, de dockerimplementation.DockerImplementation) bool {
	dei, err := de.Info(ctx)
	if err != nil {
		return false
	}
	return dei.Swarm.ControlAvailable
}
