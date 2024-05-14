package exec

import (
	"context"
	"io"
	"io/fs"

	"github.com/Mirantis/launchpad/pkg/host"
)

const (
	HostRoleExecutor = "execute"
)

func HostGetExecutor(h *host.Host) HostExecutor {
	hge := h.MatchPlugin(HostRoleExecutor)
	if hge == nil {
		return nil
	}

	he, ok := hge.(HostExecutor)
	if !ok {
		return nil
	}

	return he
}

type HostExecutor interface {
	// Connect to the host to confirm that a connection is possible
	Connect(ctx context.Context) error

	// Exec execute a command
	Exec(ctx context.Context, cmd string, inr io.Reader, opts ExecOptions) (string, string, error)

	// InstallPackages install the passed system packages
	InstallPackages(ctx context.Context, packages []string) error
	// RemovePackages uninstall some packages
	RemovePackages(ctx context.Context, packages []string) error

	// Upload content from a io.Reader to a path on the Host
	Upload(ctx context.Context, src io.Reader, dst string, fm fs.FileMode, opts ExecOptions) error

	// ServiceEnable activate and start system services
	ServiceEnable(ctx context.Context, services []string) error
	// ServiceRestart stop and restart system services
	ServiceRestart(ctx context.Context, services []string) error
	// ServiceDisable stop and disable system services
	ServiceDisable(ctx context.Context, services []string) error
}
