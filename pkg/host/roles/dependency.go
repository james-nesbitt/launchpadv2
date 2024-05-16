package roles

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
)

/**
 * Register Dependency Fullfilling for Roles Requirements
 *
 * Host components will then use this code when deciding if
 * depenendencies can be fullfilled.
 */

func init() {

	host.RegisterPluginDependencyProvider("roles", func(ctx context.Context, hc *host.HostsComponent, req dependency.Requirement) (dependency.Dependency, error) {
		if rhr, ok := req.(RequiresHostsRoles); ok {
			f := rhr.RequireHostsRoles(ctx)

			d := host.NewHostsDependency(
				fmt.Sprintf("%s:%s:%s", hc.Name(), RequiresHostsRolesType, req.Id()),
				req.Describe(),
				HostDependencyFilterFactory(hc.Hosts(), f),
			)
			slog.DebugContext(ctx, fmt.Sprintf("%s host plugin added a HostDependency '%s'for Requirement '%s'", "roles", d.Id(), req.Id()), slog.Any("dependency", d), slog.Any("requirement", req))

			return d, nil
		}
		return nil, nil
	})

}

// HostDependencyFilterFactory factory function for the Host Dependency for filtering on roles
func HostDependencyFilterFactory(hs host.Hosts, filter HostsRolesFilter) func(context.Context) (host.Hosts, error) {
	return func(_ context.Context) (host.Hosts, error) {
		return FilterRoles(hs, filter)
	}
}
