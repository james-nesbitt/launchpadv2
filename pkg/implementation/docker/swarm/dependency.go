package swarm

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

type DockerSwarmDependency interface {
	ProvidesDockerSwarm(ctx context.Context) *Swarm
}

func NewDockerSwarmDependency(id string, description string, factory func(context.Context) (*Swarm, error)) *dockerSwarmDep {
	return &dockerSwarmDep{
		id:      id,
		desc:    description,
		factory: factory}
}

type dockerSwarmDep struct {
	id   string
	desc string

	factory func(context.Context) (*Swarm, error)
}

func (dhd dockerSwarmDep) Id() string {
	return dhd.id
}

func (dhd dockerSwarmDep) Describe() string {
	return dhd.desc
}

func (dhd dockerSwarmDep) Validate(context.Context) error {
	if dhd.factory == nil {
		return dependency.ErrShouldHaveHandled
	}

	return nil
}

func (dhd dockerSwarmDep) Met(ctx context.Context) error {
	_, err := dhd.factory(ctx)
	return err
}

func (dhd dockerSwarmDep) ProvidesDockerSwarm(ctx context.Context) *Swarm {
	d, _ := dhd.factory(ctx)
	return d
}
