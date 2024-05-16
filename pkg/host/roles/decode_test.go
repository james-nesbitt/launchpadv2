package roles_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/roles"
	"gopkg.in/yaml.v3"
)

func Test_RolesDecode(t *testing.T) {
	ctx := context.Background()
	h := host.NewHost("test")
	y := `
- one
- two
- three
`
	var yn yaml.Node
	if err := yaml.Unmarshal([]byte(y), &yn); err != nil {
		t.Fatalf("error occured unmarshaling: %s", err.Error())
	}

	hp, err := host.HostPluginDecoders[roles.HostRoleRole](ctx, h, yn.Decode)
	if err != nil {
		t.Fatalf("error occured building plugin: %s", err.Error())
	}

	if !hp.RoleMatch(roles.HostRoleRole) {
		t.Errorf("host plugin doesn't match the expected")
	}

	he, ok := hp.(roles.HostRoleHandler)
	if !ok {
		t.Fatal("plugin can't handle roles")
	}

	if !he.HasRole("two") {
		t.Errorf("roles plugin didn't think it had an expected role")
	}
	if he.HasRole("four") {
		t.Errorf("roles plugin thinks it had an unexpected role")
	}

}
