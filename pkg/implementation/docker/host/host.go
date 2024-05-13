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

	d *dockerimplementation.DockerExec
}

func (h *Host) Docker(ctx context.Context) *dockerimplementation.DockerExec {
	if h.d == nil {
		h.Connect(ctx)
		def := func(ctx context.Context, cmd string, i io.Reader, rops dockerimplementation.RunOptions) (string, string, error) {
			hopts := host.ExecOptions{Sudo: h.SudoDocker}

			if rops.ShowOutput {
				hopts.OutputLevel = "info"
			}
			if rops.ShowError {
				hopts.ErrorLevel = "warn"
			}

			return h.Exec(ctx, cmd, i, hopts)
		}

		h.d = dockerimplementation.NewDockerExec(def)
	}
	return h.d
}

func (h *Host) IsSwarmManager(ctx context.Context) bool {
	i, ierr := h.Docker(ctx).Info(ctx)
	if ierr != nil {
		return false
	}

	return i.Swarm.ControlAvailable
}
