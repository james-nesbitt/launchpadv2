package mcr

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

/**
 * MCR has a HostPlugin which allows MCR specific data to be included in the
 * Host block, and also provided a Host specific docker implementation.
 */

const (
	HostRoleMCR = "mcr"
)

func init() {
	host.RegisterPluginDecoder(HostRoleMCR, func(ctx context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
		p := mcrHost{h: h}

		err := d(&p)

		return &p, err
	})
}

func HostGetMCR(h *host.Host) HostMCR {
	hgmcr := h.MatchPlugin(HostRoleMCR)
	if hgmcr == nil {
		return nil
	}

	hmcr, ok := hgmcr.(HostMCR)
	if !ok {
		return nil
	}

	return hmcr
}

type HostMCR interface {
	// MCRConfig get the daemon.json contents for the host
	MCRConfig() mcrHost
	// SudoDocker should docker be run with sudo on this host
	SudoDocker() bool
}

type mcrHost struct {
	h                *host.Host
	ShouldSudoDocker bool   `yaml:"sudo_docker"`
	DaemonJson       string `yaml:"daemon_json"`
}

func (mhc mcrHost) Id() string {
	return "mcr"
}

func (mhc mcrHost) Validate() error {
	return nil
}

func (mhc mcrHost) RoleMatch(role string) bool {
	switch role {
	case HostRoleMCR:
		return true
	case dockerhost.HostRoleDocker:
		return true
	}

	return false
}

func (mhc mcrHost) MCRConfig() string {
	return mhc.DaemonJson
}

func (mhc mcrHost) SudoDocker() bool {
	return mhc.ShouldSudoDocker
}
