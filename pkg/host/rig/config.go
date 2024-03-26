package righost

import "github.com/k0sproject/rig/v2"

// HostConfig configuration for an individual host
type Config struct {
	rig.ClientWithConfig `yaml:",inline"`

	Role string `yaml:"role"`
}
