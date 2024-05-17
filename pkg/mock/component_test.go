package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
)

func Test_MockComponentSanity(t *testing.T) {
	debug := "test"
	validate := errors.New("an error")

	var mc component.Component = mock.NewComponentWDependencies(
		"one",
		debug,
		validate,
		dependency.Requirements{
			mock.Requirement("first", "handled as the first", nil),
		},
		mock.Dependency("only", "should get used for all requirements", nil, nil),
		dependency.ErrNotHandled,
	)

	if mc.Name() != "one" {
		t.Errorf("Unexpected component name: %s", mc.Name())
	}

	if err := mc.Validate(context.Background()); !errors.Is(err, validate) {
		t.Errorf("wrong errors returned: %s", err.Error())
	}
}

func Test_MockComponentDependencySanity(t *testing.T) {
	ctx := context.Background()
	mc := mock.NewComponentWDependencies(
		"one",
		nil,
		nil,
		dependency.Requirements{
			mock.Requirement("first", "handled as the first", nil),
		},
		mock.Dependency("only", "should get used for all requirements", nil, nil),
		dependency.ErrNotHandled,
	)

	var mcr dependency.RequiresDependencies = mc

	rs := mcr.RequiresDependencies(ctx)

	if len(rs) == 0 {
		t.Errorf("Mock component didn't return any requirements: %+v", mc)
	}
	r := rs[0]
	if r.Describe() != "handled as the first" {
		t.Errorf("Mock component returned wrong requirement: %+v", mc)
	}

	var mcd dependency.ProvidesDependencies = mc

	d, err := mcd.ProvidesDependencies(ctx, nil)
	if err == nil {
		t.Errorf("Mock returned an unexpected error: %+v", mc)
	}

	if d == nil {
		t.Errorf("Mock returned a nil dependency: %+v", mc)
	} else if d.Id() != "only" {
		t.Errorf("Mock returned the wrong dependeny: %+v", mc)
	}
}
