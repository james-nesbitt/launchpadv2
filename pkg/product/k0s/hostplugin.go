package k0s

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

/**
 * K0S has a HostPlugin which allows K0S specific data to be included in the
 * Host block, and also provided a Host specific docker implementation.
 */

const (
	HostRoleK0S = "k0s"
)

var (
	// K0SManagerHostRoles the Host roles accepted for managers.
	K0SManagerHostRoles = []string{"manager"}
)

func init() {
	host.RegisterHostPluginFactory(HostRoleK0S, &k0sHostPluginFactory{})
}

type k0sHostPluginFactory struct {
	ps []*k0sHostPlugin
}

// HostPlugin build a new host plugin
func (mpf *k0sHostPluginFactory) HostPlugin(_ context.Context, h *host.Host) host.HostPlugin {
	p := &k0sHostPlugin{
		h: h,
	}
	mpf.ps = append(mpf.ps, p)

	return p
}

// Decoder provide a Host Plugin decoder function
//
// The decoder function is ugly, but it is meant to to take a
// yaml/json .HostPluginDecode() function, and turn it into a plugin
func (mpf *k0sHostPluginFactory) HostPluginDecode(_ context.Context, h *host.Host, d func(interface{}) error) (host.HostPlugin, error) {
	p := &k0sHostPlugin{h: h}
	mpf.ps = append(mpf.ps, p)

	err := d(p)

	return p, err
}

// Get the K0S plugin from a Host
func HostGetK0S(h *host.Host) *k0sHostPlugin {
	hgk0s := h.MatchPlugin(HostRoleK0S)
	if hgk0s == nil {
		return nil
	}

	hk0s, ok := hgk0s.(*k0sHostPlugin)
	if !ok {
		return nil
	}

	return hk0s
}

// k0sHostPlugin
//
// Implements no generic plugin interfaces.
type k0sHostPlugin struct {
	h       *host.Host
	K0sRole string `yaml:"role"`
	Config  string `yaml:"config"`
}

func (mhc k0sHostPlugin) Id() string {
	return "k0s"
}

func (mhc k0sHostPlugin) Validate() error {
	return nil
}

func (mhc k0sHostPlugin) RoleMatch(role string) bool {
	switch role {
	case HostRoleK0S:
		return true
	}

	return false
}

func (mhc k0sHostPlugin) K0SConfig() string {
	return mhc.Config
}

func (mhc k0sHostPlugin) IsController() bool {
	for _, r := range K0SManagerHostRoles {
		if mhc.K0sRole == r {
			return true
		}
	}
	return false
}
