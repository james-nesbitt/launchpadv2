package mock

import (
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

func NewMockPlugin(h *host.Host) mockPlugin {
	return mockPlugin{
		h: h,
	}
}

type mockPlugin struct {
	h *host.Host

	Network network.Network `yaml:"network"`
}

func (p mockPlugin) Id() string {
	return "mock"
}

func (p mockPlugin) Validate() error {
	return nil
}

func (p mockPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleMock:
		return true
	}
	return false
}
