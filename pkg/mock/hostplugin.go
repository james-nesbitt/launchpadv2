package mock

import (
	"context"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

const (
	HostRoleMock = "mock"
)

func init() {
	host.RegisterHostPluginFactory(HostRoleMock, &MockHostPluginFactory{})
}

type MockHostPluginFactory struct {
	ps []*mockHostPlugin
}

// HostPlugin build a new host plugin.
func (mpf *MockHostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &mockHostPlugin{
		h: h,
		n: &network.Network{},
	}
	mpf.ps = append(mpf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin.
func (mpf *MockHostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	hp := &mockHostPlugin{
		h: h,
		n: &network.Network{},
	}
	mpf.ps = append(mpf.ps, hp)

	if err := d(hp); err != nil {
		return hp, err
	}
	return hp, nil
}

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

func NewMockHostPlugin(h *host.Host, n *network.Network) mockHostPlugin {
	return mockHostPlugin{
		h: h,
		n: n,
	}
}

type mockHostPlugin struct {
	h *host.Host
	n *network.Network `yaml:"network"`
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
