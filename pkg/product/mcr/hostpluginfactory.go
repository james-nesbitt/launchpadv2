package mcr

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

func init() {
	host.RegisterHostPluginFactory(HostRoleMCR, &hostPluginFactory{})
}

type hostPluginFactory struct {
	ps []*hostPlugin
}

// HostPlugin build a new host plugin.
func (pf *hostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &hostPlugin{
		h: h,
	}
	pf.ps = append(pf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin.
func (pf *hostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &hostPlugin{h: h}

	if err := d(p); err != nil {
		return p, err
	}

	pf.ps = append(pf.ps, p)

	return p, nil
}
