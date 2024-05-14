package roles

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

// HostDependencyFilterFactory factory function for the Host Dependency for filtering on roles
func HostDependencyFilterFactory(hs host.Hosts, filter HostsRolesFilter) func(context.Context) (host.Hosts, error) {
	return func(_ context.Context) (host.Hosts, error) {
		return FilterRoles(hs, filter)
	}
}
