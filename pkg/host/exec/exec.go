package exec

import (
	"context"
	"io"

	"github.com/Mirantis/launchpad/pkg/host"
)

// HostExecutor
//
// A plugin role that indicated that commands can be
// executed on the host.
// Additionally the host can manage services and packages.
//
// Currently we have only one plugin for this role: rig

const (
	HostRoleExecutor = "execute"
)

// HostGetExecutor get the host plugin that can execute.
//
// If the hosts doesn't have an appropriate plugin then
// nil is returned.
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

	// ExecInteractive get a shell
	ExecInteractive(ctx context.Context, opts ExecOptions) error

	// InstallPackages install the passed system packages
	InstallPackages(ctx context.Context, packages []string) error
	// RemovePackages uninstall some packages
	RemovePackages(ctx context.Context, packages []string) error

	// ServiceEnable activate and start system services
	ServiceEnable(ctx context.Context, services []string) error
	// ServiceRestart stop and restart system services
	ServiceRestart(ctx context.Context, services []string) error
	// ServiceDisable stop and disable system services
	ServiceDisable(ctx context.Context, services []string) error
	// ServiceIsRunning is the service running
	ServiceIsRunning(ctx context.Context, services []string) error
}
