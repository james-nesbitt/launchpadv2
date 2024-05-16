package rig

import (
	"github.com/k0sproject/rig/v2"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

type rigHostPlugin struct {
	h   *host.Host
	rig *rig.ClientWithConfig
}

func (mhc rigHostPlugin) Id() string {
	return "rig"
}

func (mhc rigHostPlugin) Validate() error {
	return nil
}

func (mhc rigHostPlugin) RoleMatch(role string) bool {
	switch role {
	case host.HostRoleDiscover:
		return true
	case exec.HostRoleExecutor:
		return true
	case network.HostRoleNetwork:
		return true
	case exec.HostRolePlatform:
		return true
	}
	return false
}

func (p rigHostPlugin) hid() string {
	return p.h.Id()
}
