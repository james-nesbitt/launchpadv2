package roles_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/roles"
)

func Test_RolesPlugin_PluginSanty(t *testing.T) {
	h := host.NewHost("test")
	var rp host.HostPlugin = roles.NewHostRolesPlugin([]string{"test"})

	h.AddPlugin(rp)

	hrp := h.MatchPlugin(roles.HostRoleRole)
	if hrp == nil {
		t.Error("host couldn't match roles plugin")
	}
}

func Test_RolesPlugin_MatchRoles(t *testing.T) {
	rs := []string{"one", "two", "three"}
	rp := roles.NewHostRolesPlugin(rs)

	if !rp.HasRole("one") {
		t.Errorf("plugin didn't match expected role: %s", "one")
	}
	if rp.HasRole("four") {
		t.Errorf("plugin matched unexpected role: %s", "four")
	}

}
