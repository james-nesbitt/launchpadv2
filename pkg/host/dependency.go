package host

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// Provides a dependency for some type of Requirements.
func (hc *HostsComponent) ProvidesDependencies(ctx context.Context, req dependency.Requirement) (dependency.Dependency, error) {
	if fhr, ok := req.(RequiresFilteredHosts); ok {
		// RequiersFiltertHosts requirements want a dependency that will return a filtered set of the component hosts.
		// We use the HostDependency to fulfill that, and use the FilterFactory function as the dependency factory.

		ff := fhr.HostsFilter(ctx) // get the filter object from the requirement

		d := NewHostsDependency(
			fmt.Sprintf("%s:%s:%s", hc.id, "filter", req.Id()),
			req.Describe(),
			HostsDependencyFilterFactory(hc.hosts, ff),
		)
		slog.DebugContext(ctx, fmt.Sprintf("%s added a HostDependency '%s'for FilteredHosts Requirement '%s'", ComponentType, d.Id(), req.Id()), slog.Any("dependency", d), slog.Any("requirement", req))

		hc.deps = append(hc.deps, d)

		return d, nil
	}

	return nil, nil
}
