package host

import (
	"github.com/k0sproject/rig"

	"github.com/Mirantis/launchpad/pkg/product/common"
)

// HostConfig configuration for an individual host
type HostConfig struct {
	rig.Connection `yaml:",inline"`

	Role             string            `yaml:"role" validate:"oneof=manager worker msr"`
	PrivateInterface string            `yaml:"privateInterface,omitempty" validate:"omitempty,gt=2"`
	Environment      map[string]string `yaml:"environment,flow,omitempty" default:"{}"`
	Hooks            common.Hooks      `yaml:"hooks,omitempty" validate:"dive,keys,oneof=apply reset,endkeys,dive,keys,oneof=before after,endkeys,omitempty"`
	ImageDir         string            `yaml:"imageDir,omitempty"`
}
