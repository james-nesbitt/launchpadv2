package dockerhost

import (
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
)

func NewDockerHosts(hosts host.Hosts, version docker.Version) *DockerHosts {
	return &DockerHosts{
		hosts:   hosts,
		version: version,
	}
}

type DockerHosts struct {
	hosts   host.Hosts
	version docker.Version
}
