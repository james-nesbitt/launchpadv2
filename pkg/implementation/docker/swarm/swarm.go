package swarm

import (
	"fmt"

	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
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

// SwarmAddress convert an IP into a address
func SwarmAddress(ip string) string {
	return fmt.Sprintf("%s:%d", ip, 2377)
}
