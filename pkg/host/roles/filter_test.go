package roles_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/roles"
)

func Test_RolesFilter(t *testing.T) {
	one := host.NewHost("one")
	one.AddPlugin(roles.NewHostRolesPlugin([]string{"A", "B"}))
	two := host.NewHost("two")
	three := host.NewHost("three")
	three.AddPlugin(roles.NewHostRolesPlugin([]string{"B", "C"}))

	hs := host.NewHosts(
		one,
		two,
		three,
	)

	if fhs, err := roles.FilterRoles(hs, roles.HostsRolesFilter{Roles: []string{"B"}, Min: 2}); err != nil {
		t.Errorf("filter roles failed: %s", err.Error())
	} else if len(fhs) != 2 {
		t.Errorf("wrong number of results from filter")
	}

	if fhs, err := roles.FilterRoles(hs, roles.HostsRolesFilter{Roles: []string{"C"}}); err != nil {
		t.Errorf("filter roles failed: %s", err.Error())
	} else if len(fhs) != 1 {
		t.Errorf("wrong number of results from filter: %+v", fhs)
	}

	if fhs, err := roles.FilterRoles(hs, roles.HostsRolesFilter{Roles: []string{"D"}}); err != nil {
		t.Errorf("filter roles failed: %s", err.Error())
	} else if len(fhs) != 0 {
		t.Errorf("wrong number of results from filter: %+v", fhs)
	}
}
