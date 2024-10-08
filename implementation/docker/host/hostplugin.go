package dockerhost

import (
	dockerimplementation "github.com/Mirantis/launchpad/implementation/docker"
	"github.com/Mirantis/launchpad/pkg/host"
)

var (
	// HostRoleDocker indicates that a Host plugin can produce a Docker client.
	HostRoleDockerExec = "docker-exec"
)

// HostGetDockerExec get a DockerExec from a host.
func HostGetDockerExec(h *host.Host) *dockerimplementation.DockerExec {
	hgde := h.MatchPlugin(HostRoleDockerExec)
	if hgde == nil {
		return nil
	}

	hde, ok := hgde.(HostDockerExec)
	if !ok {
		return nil
	}

	return hde.DockerExec()
}

type HostDockerExec interface {
	DockerExec() *dockerimplementation.DockerExec
}
