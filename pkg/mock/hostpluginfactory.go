package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

func init() {
	host.RegisterHostPluginFactory(HostRoleMock, &MockHostPluginFactory{})
}

// MockHostPluginFactory a host plugin factory which produces mock host plugins.
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
