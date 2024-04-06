package dockerhost

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
)

const (
	ImplementationType = "docker-hosts"
)

type DockerHostsRequirement interface {
	NeedsDockerHost(context.Context) docker.Version
}

func NewDockerHostsRequirement(id string, desc string, version docker.Version) *dockerHostsReq {
	return &dockerHostsReq{
		id:   id,
		desc: desc,
		ver:  version,
	}
}

type DockerHostsDependency interface {
	ProvidesDockerHost(ctx context.Context) *DockerHosts
}

func NewDockerHostsDependency(id string, description string, factory func(context.Context) (*DockerHosts, error)) *dockerHostsDep {
	return &dockerHostsDep{
		id:      id,
		desc:    description,
		factory: factory}
}

type dockerHostsReq struct {
	id   string
	desc string

	ver docker.Version

	dep dependency.Dependency
}

func (dhr dockerHostsReq) Id() string {
	return dhr.id
}

func (dhr dockerHostsReq) Describe() string {
	return dhr.desc
}

func (dhr *dockerHostsReq) Match(dep dependency.Dependency) error {
	if _, ok := dep.(DockerHostsDependency); !ok {
		return dependency.ErrDependencyNotMatched
	}

	dhr.dep = dep
	return nil
}

func (dhr dockerHostsReq) Matched(context.Context) dependency.Dependency {
	return dhr.dep
}

func (dhr dockerHostsReq) NeedsDockerHost(_ context.Context) docker.Version {
	return dhr.ver
}

type dockerHostsDep struct {
	id   string
	desc string

	factory func(context.Context) (*DockerHosts, error)
}

func (dhd dockerHostsDep) Id() string {
	return dhd.id
}

func (dhd dockerHostsDep) Describe() string {
	return dhd.desc
}

func (dhd dockerHostsDep) Validate(context.Context) error {
	if dhd.factory == nil {
		return dependency.ErrDependencyShouldHaveHandled
	}

	return nil
}

func (dhd dockerHostsDep) Met(ctx context.Context) error {
	_, err := dhd.factory(ctx)
	return err
}

func (dhd dockerHostsDep) ProvidesDockerHost(ctx context.Context) *DockerHosts {
	d, _ := dhd.factory(ctx)
	return d
}
