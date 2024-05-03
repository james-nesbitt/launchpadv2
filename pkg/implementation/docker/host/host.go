package dockerhost

import (
	"context"
	"io"

	"github.com/Mirantis/launchpad/pkg/host"
	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
)

// Host a host.Host that can retrieve a DockerExec implementation.
type Host struct {
	host.Host

	SudoDocker bool
}

func (h Host) Docker(ctx context.Context) dockerimplementation.DockerExec {
	def := func(ctx context.Context, cmd string, i io.Reader) (string, string, error) {
		return h.Exec(ctx, cmd, i, host.ExecOptions{})
	}
	if h.SudoDocker {
		def = func(ctx context.Context, cmd string, i io.Reader) (string, string, error) {
			return h.Exec(ctx, cmd, i, host.ExecOptions{Sudo: true})
		}
	}

	return dockerimplementation.NewDockerExec(def) // use the host exec as a DockerExec executor
}
