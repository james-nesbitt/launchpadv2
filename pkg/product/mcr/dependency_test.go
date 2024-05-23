package mcr_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/product/mcr"
)

func Test_MCRDependencySanity(t *testing.T) {
	mcrcl := mcr.NewComponent(
		"dummy",
		mcr.Config{},
	)

	if _, ok := interface{}(mcrcl).(dependency.RequiresDependencies); !ok {
		t.Errorf("MCR fails to advertise dependencies")
	}

	if _, ok := interface{}(mcrcl).(dependency.ProvidesDependencies); !ok {
		t.Errorf("MCR fails to advertise dependency fullfillment")
	}
}
