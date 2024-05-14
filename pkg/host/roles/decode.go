package roles

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

func init() {
	host.RegisterPluginDecoder(HostRoleRole, func(ctx context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
		p := hostRoles{}

		err := d(&p)

		return p, err
	})
}
