package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

const (
	HostRoleMock = "mock"
)

func init() {
	host.RegisterPluginDecoder(HostRoleMock, func(ctx context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
		hp := mockPlugin{
			h:       h,
			Network: network.Network{},
		}

		if err := d(&hp); err != nil {
			return hp, err
		}
		return hp, nil

	})
}

func HostGetMock(h *host.Host) *mockPlugin {
	hgm := h.MatchPlugin(HostRoleMock)
	if hgm == nil {
		return nil
	}

	hm, ok := hgm.(mockPlugin)
	if !ok {
		return nil
	}

	return &hm
}
