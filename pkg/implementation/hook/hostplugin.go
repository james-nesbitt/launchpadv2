package hook

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

const (
	HostRoleHooks = "hooks"
)

// HostGetHooks get hooks for a host
//
// If the hosts doesn't have an appropriate plugin then
// nil is returned.
func HostGetExecutor(h *host.Host) Hooks {
	hge := h.MatchPlugin(HostRoleHooks)
	if hge == nil {
		return nil
	}

	he, ok := hge.(HostHooks)
	if !ok {
		return nil
	}

	return he.Hooks()
}

// HostHooks this plugin can deliver hooks
type HostHooks interface {
	// Provides hooks
	Hooks() Hooks
}

// register our host plugin factory
func init() {
	host.RegisterHostPluginFactory(HostRoleHooks, &hookHostPluginFactory{})
}

type hookHostPluginFactory struct {
	ps []*hookHostPlugin
}

// Plugin build a new host plugin
func (rpf *hookHostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &hookHostPlugin{
		h:     h,
		hooks: Hooks{},
	}

	rpf.ps = append(rpf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .Decode() function, and turn it into a plugin
func (rpf *hookHostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &hookHostPlugin{
		h:     h,
		hooks: Hooks{},
	}

	if err := d(p.hooks); err != nil {
		return p, err
	}

	rpf.ps = append(rpf.ps, p)

	return p, nil

}

// hookhHostPlugin a host plugin that uses inline hooks
//
// Implements:
// HostRoleHooks: provides hooks
type hookHostPlugin struct {
	h     *host.Host
	hooks Hooks
}

func (mhc hookHostPlugin) Id() string {
	return "hook"
}

func (mhc hookHostPlugin) Validate() error {
	return nil
}

func (mhc hookHostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleHooks:
		return true
	}
	return false
}

func (p hookHostPlugin) hid() string {
	return p.h.Id()
}
