package swarm

import (
	"context"

	dockerimplementation "github.com/Mirantis/launchpad/implementation/docker"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ImplementationType = "docker-swarm"
)

type DockerSwarmRequirement interface {
	NeedsDockerSwarm(context.Context) dockerimplementation.Version
}

func NewDockerSwarmRequirement(id string, desc string, version dockerimplementation.Version) *dockerSwarmReq {
	return &dockerSwarmReq{
		id:   id,
		desc: desc,
		ver:  version,
	}
}

type dockerSwarmReq struct {
	id   string
	desc string

	ver dockerimplementation.Version

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

func (dhr dockerSwarmReq) NeedsDockerSwarm(_ context.Context) dockerimplementation.Version {
	return dhr.ver
}
