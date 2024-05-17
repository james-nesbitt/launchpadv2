package mcr

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
	HostRoleMCR = "mcr"
)

var (
	// MCRManagerHostRoles the Host roles accepted for managers.
	MCRManagerHostRoles = []string{"manager"}
)

func init() {
	host.RegisterHostPluginFactory(HostRoleMCR, &mcrHostPluginFactory{})
}

type mcrHostPluginFactory struct {
	ps []*mcrHostPlugin
}

// Plugin build a new host plugin
func (mpf *mcrHostPluginFactory) Plugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &mcrHostPlugin{
		h: h,
	}
	mpf.ps = append(mpf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .Decode() function, and turn it into a plugin
func (mpf *mcrHostPluginFactory) Decode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &mcrHostPlugin{h: h}
	mpf.ps = append(mpf.ps, p)

	err := d(p)

	return p, err
}

// Get the MCR plugin from a Host
func HostGetMCR(h *host.Host) *mcrHostPlugin {
	hgmcr := h.MatchPlugin(HostRoleMCR)
	if hgmcr == nil {
		return nil
	}

	hmcr, ok := hgmcr.(*mcrHostPlugin)
	if !ok {
		return nil
	}

	return hmcr
}

// mcrHostPlugin
//
// Implements no generic plugin interfaces.
type mcrHostPlugin struct {
	h                *host.Host
	SwarmRole        string `yaml:"role"`
	ShouldSudoDocker bool   `yaml:"sudo_docker"`
	DaemonJson       string `yaml:"daemon_json"`
}

func (mhc mcrHostPlugin) Id() string {
	return "mcr"
}

func (mhc mcrHostPlugin) Validate() error {
	return nil
}

func (mhc mcrHostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleMCR:
		return true
	case dockerhost.HostRoleDockerExec:
		// pkg/implementation/docker/host
		return true
	}

	return false
}

// DockerExecOptions meet the dockerhost.HostDockerExec interface
func (mhc mcrHostPlugin) DockerExec() *dockerimplementation.DockerExec {
	e := exec.HostGetExecutor(mhc.h)
	if e == nil {
		return nil
	}

	def := func(ctx context.Context, cmd string, i io.Reader, rops dockerimplementation.RunOptions) (string, string, error) {
		hopts := exec.ExecOptions{Sudo: mhc.ShouldSudoDocker}

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

func (mhc mcrHostPlugin) MCRConfig() string {
	return mhc.DaemonJson
}

func (mhc mcrHostPlugin) IsManager() bool {
	for _, r := range MCRManagerHostRoles {
		if mhc.SwarmRole == r {
			return true
		}
	}
	return false
}
