package mock

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

const (
	HostRoleMock = "mock"
)

var (
	ErrMockHostNetworkMissing = errors.New("mock host plugin had no network information")
)

func HostGetMock(h *host.Host) *mockHostPlugin {
	hgm := h.MatchPlugin(HostRoleMock)
	if hgm == nil {
		return nil
	}

	hm, ok := hgm.(*mockHostPlugin)
	if !ok {
		slog.Warn("could not convert mock plugin")
		return nil
	}

	return hm
}

func NewMockHostPlugin(h *host.Host, n *network.Network) *mockHostPlugin {
	return &mockHostPlugin{
		h: h,
		n: n,
	}
}

type mockHostPlugin struct {
	h *host.Host
	n *network.Network `yaml:"network"`
}

func (p mockHostPlugin) Id() string {
	return HostRoleMock
}

func (p mockHostPlugin) Validate() error {
	return nil
}

func (p mockHostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleMock:
		return true
	case network.HostRoleNetwork:
		return p.n != nil
	}
	return false
}

// Network provide network information.
func (p *mockHostPlugin) Network(ctx context.Context) (network.Network, error) {
	if p.n == nil {
		return network.Network{}, ErrMockHostNetworkMissing
	}

	return *p.n, nil
}
