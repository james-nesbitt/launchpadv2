package k0s

import (
	"context"

	"github.com/creasty/defaults"

	"github.com/Mirantis/launchpad/pkg/host"
)

func init() {
	host.RegisterHostPluginFactory(HostRoleK0S, &HostPluginFactory{})
}

type HostPluginFactory struct {
	ps []*hostPlugin
}

// HostPlugin build a new host plugin.
func (pf *HostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	c := HostConfig{
		Environment: map[string]string{},
	}
	p := &hostPlugin{
		h: h,
		c: c,
		s: hostState{},
	}

	defaults.Set(p)

	pf.ps = append(pf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin.
func (pf *HostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	c := HostConfig{
		Environment: map[string]string{},
	}
	p := &hostPlugin{
		h: h,
		s: hostState{},
	}

	if err := d(&c); err != nil {
		return p, err
	}

	p.c = c

	pf.ps = append(pf.ps, p)

	return p, nil
}

