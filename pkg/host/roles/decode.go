package roles

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

/**
 * Register the Roles Host plugin with the Hosts component
 *
 * This means that "roles" information for Hosts can be used
 */

func init() {
	host.RegisterPluginDecoder(HostRoleRole, func(ctx context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
		p := hostRoles{}

		err := d(&p.roles)

		return p, err
	})
}
