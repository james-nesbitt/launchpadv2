package roles_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host/roles"
)

func Test_RequirementSanity(t *testing.T) {
	ctx := context.Background()
	f := roles.HostsRolesFilter{
		Min:   1,
		Max:   2,
		Roles: []string{"one", "two"},
	}
	var r dependency.Requirement = roles.NewHostsRolesRequirement("test", "my test requirement", f)

	if r.Id() != "test" {
		t.Errorf("wrong requirement id")
	}

	hrr, ok := r.(roles.RequiresHostsRoles)
	if !ok {
		t.Fatal("host roles requirement doesn't do requiring host roles")
	}

	hrrf := hrr.RequireHostsRoles(ctx)
	if hrrf.Min != f.Min {
		t.Error("wrong min")
	}
}
