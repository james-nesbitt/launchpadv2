package mock

import (
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

func HostGetMock(h *host.Host) *mockHostPlugin {
	hgm := h.MatchPlugin(HostRoleMock)
	if hgm == nil {
		return nil
	}

	hm, ok := hgm.(mockHostPlugin)
	if !ok {
		slog.Warn("could not convert mock plugin")
		return nil
	}

	return &hm
}

func NewMockHostPlugin(h *host.Host) mockHostPlugin {
	return mockHostPlugin{
		h: h,
	}
}

type mockHostPlugin struct {
	h *host.Host

	Network network.Network `yaml:"network"`
}

func (p mockHostPlugin) Id() string {
	return "mock"
}

func (p mockHostPlugin) Validate() error {
	return nil
}

func (p mockHostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleMock:
		return true
	}
	return false
}
