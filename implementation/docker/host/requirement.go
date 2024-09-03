package dockerhost

import (
	"context"

	dockerimplementation "github.com/Mirantis/launchpad/implementation/docker"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ImplementationType = "docker-hosts"
)

type DockerHostsRequirement interface {
	NeedsDockerHost(context.Context) dockerimplementation.Version
}

func NewDockerHostsRequirement(id string, desc string, version dockerimplementation.Version) *dockerHostsReq {
	return &dockerHostsReq{
		id:   id,
		desc: desc,
		ver:  version,
	}
}

type dockerHostsReq struct {
	id   string
	desc string

	ver dockerimplementation.Version

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

func (dhr dockerHostsReq) NeedsDockerHost(_ context.Context) dockerimplementation.Version {
	return dhr.ver
}
