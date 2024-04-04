package mock_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	mockcomponent "github.com/Mirantis/launchpad/pkg/component/mock"
	"github.com/Mirantis/launchpad/pkg/dependency"
	mockdependency "github.com/Mirantis/launchpad/pkg/dependency/mock"
	"github.com/Mirantis/launchpad/pkg/phase"
)

func Test_MockSanity(t *testing.T) {
	var mc component.Component = mockcomponent.NewComponentWDependencies(
		"one",
		phase.Actions{},
		dependency.Requirements{
			mockdependency.Requirement(
				"first",
				"handled as the first",
				nil,
			),
		},
		mockdependency.Dependency(
			"only",
			"should get used for all requirements",
			nil,
		),
		dependency.ErrDependencyNotHandled,
	)

	if mc.Name() != "one" {
		t.Errorf("Unexpected component name: %s", mc.Name())
	}
}

func Test_MockDependencySanity(t *testing.T) {
	ctx := context.Background()
	mc := mockcomponent.NewComponentWDependencies(
		"one",
		phase.Actions{},
		dependency.Requirements{
			mockdependency.Requirement(
				"first",
				"handled as the first",
				nil,
			),
		},
		mockdependency.Dependency(
			"only",
			"should get used for all requirements",
			nil,
		),
		dependency.ErrDependencyNotHandled,
	)

	var mcr dependency.HasDependencies = mc

	rs := mcr.Requires(ctx)

	if len(rs) == 0 {
		t.Errorf("Mock component didn't return any requirements: %+v", mc)
	}
	r := rs[0]
	if r.Describe() != "handled as the first" {
		t.Errorf("Mock component returned wrong requirement: %+v", mc)
	}

	var mcd dependency.FullfillsDependencies = mc

	d, err := mcd.Provides(ctx, nil)
	if err == nil {
		t.Errorf("Mock returned an unexpected error: %+v", mc)
	}

	if d == nil {
		t.Errorf("Mock returned a nil dependency: %+v", mc)
	} else if d.Id() != "only" {
		t.Errorf("Mock returned the wrong dependeny: %+v", mc)
	}
}
