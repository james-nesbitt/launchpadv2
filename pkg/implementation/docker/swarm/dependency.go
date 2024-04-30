package swarm

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
)

const (
	ImplementationType = "docker-swarm"
)

type DockerSwarmRequirement interface {
	NeedsDockerSwarm(context.Context) docker.Version
}

func NewDockerSwarmRequirement(id string, desc string, version docker.Version) *dockerSwarmReq {
	return &dockerSwarmReq{
		id:   id,
		desc: desc,
		ver:  version,
	}
}

type DockerSwarmDependency interface {
	ProvidesDockerSwarm(ctx context.Context) *Swarm
}

func NewDockerSwarmDependency(id string, description string, factory func(context.Context) (*Swarm, error)) *dockerSwarmDep {
	return &dockerSwarmDep{
		id:      id,
		desc:    description,
		factory: factory}
}

type dockerSwarmReq struct {
	id   string
	desc string

	ver docker.Version

	dep dependency.Dependency
}

func (dhr dockerSwarmReq) Id() string {
	return dhr.id
}

func (dhr dockerSwarmReq) Describe() string {
	return dhr.desc
}

func (dhr *dockerSwarmReq) Match(dep dependency.Dependency) error {
	if _, ok := dep.(DockerSwarmDependency); !ok {
		return dependency.ErrDependencyNotMatched
	}

	dhr.dep = dep
	return nil
}

func (dhr dockerSwarmReq) Matched(context.Context) dependency.Dependency {
	return dhr.dep
}

func (dhr dockerSwarmReq) NeedsDockerSwarm(_ context.Context) docker.Version {
	return dhr.ver
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
