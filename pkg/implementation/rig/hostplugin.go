package rig

import (
	"context"

	"github.com/k0sproject/rig/v2"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

const (
	HostRoleRig = "rig"
)

func init() {
	host.RegisterHostPluginFactory(HostRoleRig, &RigHostPluginFactory{})
}

type RigHostPluginFactory struct {
	ps []*rigHostPlugin
}

// HostPlugin build a new host plugin
func (rpf *RigHostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &rigHostPlugin{
		h:   h,
		rig: &rig.ClientWithConfig{},
	}

	rpf.ps = append(rpf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin
func (rpf *RigHostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &rigHostPlugin{
		rig: &rig.ClientWithConfig{},
		h:   h,
	}

	if err := d(p.rig); err != nil {
		return p, err
	}

	rpf.ps = append(rpf.ps, p)

	return p, nil
}

// righHostPlugin a host plugin that uses rig a a backend
//
// Implements:
// host.HostRoleDiscover: uses a rig.Client.Connect()
// exec.HostRoleExecutor: uses all of the rig.Client functions
// network.HostRoleNetwork: uses some internal functions to discover network
// exec.HostRolePlatform: uses some internal functions to discover platform
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
	case exec.HostRoleFiles:
		return true
	case exec.HostRolePlatform:
		return true
	}
	return false
}

func (p rigHostPlugin) hid() string {
	return p.h.Id()
}
