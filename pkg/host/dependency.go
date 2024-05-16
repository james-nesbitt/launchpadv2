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
		ff := fhr.RequireFilteredHosts(ctx)

		d := NewHostsDependency(
			fmt.Sprintf("%s:%s:%s", hc.id, "filter", req.Id()),
			req.Describe(),
			func(ctx context.Context) (Hosts, error) {
				return ff(ctx, hc.hosts)
			},
		)
		slog.DebugContext(ctx, fmt.Sprintf("%s added a HostDependency '%s'for Requirement '%s'", ComponentType, d.Id(), req.Id()), slog.Any("dependency", d), slog.Any("requirement", req))

		hc.deps = append(hc.deps, d)

		return d, nil
	}

	return nil, nil
}
