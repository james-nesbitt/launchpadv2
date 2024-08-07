package dockerhost

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
)

// DockerHostsDependency a Dependency that can produce hosts that have a Docker client.
type DockerHostsDependency interface {
	// ProvidesDockerHost provide hosts that can produce a Docker client
	ProvidesDockerHost(context.Context) host.Hosts
}

// NewDockerHostsDependency produce a new Dependency that can produce hosts that have Docker clients from a factory.
func NewDockerHostsDependency(id string, description string, factory func(context.Context) (host.Hosts, error)) *dockerHostsDep {
	return &dockerHostsDep{
		id:      id,
		desc:    description,
		factory: factory,
	}
}

type dockerHostsDep struct {
	id   string
	desc string

	factory func(context.Context) (host.Hosts, error)

	events dependency.Events
}

func (dhd dockerHostsDep) Id() string {
	return dhd.id
}

func (dhd dockerHostsDep) Describe() string {
	return dhd.desc
}

func (dhd dockerHostsDep) Validate(context.Context) error {
	if dhd.factory == nil {
		return dependency.ErrShouldHaveHandled
	}

	return nil
}

func (dhd dockerHostsDep) Met(ctx context.Context) error {
	_, err := dhd.factory(ctx)
	return err
}

func (dhd dockerHostsDep) ProvidesDockerHost(ctx context.Context) host.Hosts {
	d, _ := dhd.factory(ctx)
	return d
}

func (dhd dockerHostsDep) DeliversEvents(context.Context) dependency.Events {
	if dhd.events == nil {
		dhd.events = dependency.Events{
			dependency.EventKeyActivated: &dependency.Event{
				Id: fmt.Sprintf("%s:%s", dhd.Id(), dependency.EventKeyActivated),
			},
			dependency.EventKeyDeActivated: &dependency.Event{
				Id: fmt.Sprintf("%s:%s", dhd.Id(), dependency.EventKeyDeActivated),
			},
		}
	}

	return dhd.events
}
