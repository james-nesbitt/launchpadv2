package mkex

import (
	"context"
	"io"

	dockerimplementation "github.com/Mirantis/launchpad/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

/**
 * MCR has a HostPlugin which allows MCR specific data to be included in the
 * Host block, and also provided a Host specific docker implementation.
 */

const (
	HostRoleMKEX = "mkex"
)

var (
	// ManagerHostRoles the Host roles accepted for managers.
	ManagerHostRoles = []string{"manager"}
)

// Get the MKEX plugin from a Host.
func HostGetMKEX(h *host.Host) *hostPlugin {
	hgmcr := h.MatchPlugin(HostRoleMKEX)
	if hgmcr == nil {
		return nil
	}

	hmcr, ok := hgmcr.(*hostPlugin)
	if !ok {
		return nil
	}

	return hmcr
}

// hostPlugin
//
// Implements no generic plugin interfaces.
type hostPlugin struct {
	h          *host.Host
	SwarmRole  string `yaml:"role"`
	DaemonJson string `yaml:"daemon_json"`
}

func (_ hostPlugin) Id() string {
	return "mcr"
}

func (_ hostPlugin) Validate() error {
	return nil
}

func (hc hostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleMKEX:
		return true
	case dockerhost.HostRoleDockerExec:
		// pkg/implementation/docker/host
		return true
	}

	return false
}

// DockerExecOptions meet the dockerhost.HostDockerExec interface.
func (hc hostPlugin) DockerExec() *dockerimplementation.DockerExec {
	e := exec.HostGetExecutor(hc.h)
	if e == nil {
		return nil
	}

	def := func(ctx context.Context, cmd string, i io.Reader, rops dockerimplementation.RunOptions) (string, string, error) {
		hopts := exec.ExecOptions{Sudo: false}

		if rops.ShowOutput {
			hopts.OutputLevel = "info"
		}
		if rops.ShowError {
			hopts.ErrorLevel = "warn"
		}

		return e.Exec(ctx, cmd, i, hopts)
	}

	return dockerimplementation.NewDockerExec(def)
}

func (hc hostPlugin) DockerConfig() string {
	return hc.DaemonJson
}

func (hc hostPlugin) IsManager() bool {
	for _, r := range ManagerHostRoles {
		if hc.SwarmRole == r {
			return true
		}
	}
	return false
}
