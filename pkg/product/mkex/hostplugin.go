package mkex

import (
	"context"
	"io"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	dockerimplementation "github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
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

func init() {
	host.RegisterHostPluginFactory(HostRoleMKEX, &hostPluginFactory{})
}

type hostPluginFactory struct {
	ps []*hostPlugin
}

// HostPlugin build a new host plugin
func (mpf *hostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &hostPlugin{
		h: h,
	}
	mpf.ps = append(mpf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin
func (mpf *hostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &hostPlugin{h: h}
	mpf.ps = append(mpf.ps, p)

	err := d(p)

	return p, err
}

// Get the MKEX plugin from a Host
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

// DockerExecOptions meet the dockerhost.HostDockerExec interface
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
