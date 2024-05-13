package host

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ComponentType = "hosts"
)

func NewHostsComponent(id string, hosts Hosts) *HostsComponent {
	return &HostsComponent{
		id:    id,
		hosts: hosts,
		deps:  dependency.Dependencies{},
	}
}

type HostsComponent struct {
	id    string
	hosts Hosts

	deps dependency.Dependencies
}

func (hc HostsComponent) Name() string {
	if hc.id == ComponentType {
		return hc.id
	}
	return fmt.Sprintf("%s:%s", ComponentType, hc.id)
}

func (hc HostsComponent) Debug() interface{} {
	return struct {
		ID    string
		Hosts Hosts
	}{
		ID:    hc.id,
		Hosts: hc.hosts,
	}
}

func (hc HostsComponent) Validate(_ context.Context) error {
	return nil
}

// Provides a dependency for some type of Requirements.
func (hc *HostsComponent) ProvidesDependencies(ctx context.Context, req dependency.Requirement) (dependency.Dependency, error) {
	if rhr, ok := req.(RequiresHostsRoles); ok {
		f := rhr.RequireHostsRoles(ctx)

		d := NewHostsDependency(
			fmt.Sprintf("%s:%s:%s", hc.id, RequiresHostsRolesType, req.Id()),
			req.Describe(),
			func(ctx context.Context) (Hosts, error) {
				hs, err := hc.hosts.FilterRoles(f)
				return hs, err
			},
		)
		slog.DebugContext(ctx, fmt.Sprintf("%s added a HostDependency '%s'for Requirement '%s'", ComponentType, d.Id(), req.Id()), slog.Any("dependency", d), slog.Any("requirement", req))

		hc.deps = append(hc.deps, d)

		return d, nil
	}

	return nil, nil
}
