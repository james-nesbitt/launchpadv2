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
	host.RegisterHostPluginFactory(HostRoleMock, &MockHostPluginFactory{})
}

type MockHostPluginFactory struct {
	ps []*mockHostPlugin
}

// Plugin build a new host plugin
func (mpf *MockHostPluginFactory) Plugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &mockHostPlugin{
		h:       h,
		Network: network.Network{},
	}
	mpf.ps = append(mpf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .Decode() function, and turn it into a plugin
func (mpf *MockHostPluginFactory) Decode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	hp := &mockHostPlugin{
		h:       h,
		Network: network.Network{},
	}
	mpf.ps = append(mpf.ps, hp)

	if err := d(hp); err != nil {
		return hp, err
	}
	return hp, nil
}
