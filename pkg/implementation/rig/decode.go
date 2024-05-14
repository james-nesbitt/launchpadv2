package rig

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/k0sproject/rig/v2"
)

const (
	HostRoleRig = "rig"
)

var (
	sbinPath = "PATH=/usr/local/sbin:/usr/sbin:/sbin:$PATH"
)

func init() {
	host.RegisterPluginDecoder(HostRoleRig, func(ctx context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
		hp := rigHostPlugin{
			rig: &rig.ClientWithConfig{},
			h:   h,
		}

		if err := d(hp.rig); err != nil {
			return hp, err
		}
		return hp, nil

	})
}
