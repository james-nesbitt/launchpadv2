package project_test

import (
	"context"
	"errors"
	"testing"

	mockcomponent "github.com/Mirantis/launchpad/pkg/component/mock"
	"github.com/Mirantis/launchpad/pkg/dependency"
	mockdependency "github.com/Mirantis/launchpad/pkg/dependency/mock"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/project"
)

func Test_ProjectComponentValidation(t *testing.T) {
	ctx := context.Background()
	e := errors.New("validation error")
	cp := mockcomponent.NewComponent(
		"mock",
		nil,
		e,
	)

	if err := cp.Validate(ctx); err == nil {
		t.Errorf("Mock component was supposed to fail validation : %+v", cp)
	}

	cl := project.Project{
		Components: component.Components{
			"mock": cp,
		},
	}

	if err := cl.Validate(ctx); err == nil {
		t.Errorf("expected error, none returned in validation: %+v", cl)
	}

}

func Test_ProjectDependencyValidation(t *testing.T) {
	ctx := context.Background()
	cl := project.Project{
		Components: component.Components{
			"one": mockcomponent.NewComponentWDependencies(
				"one",
				nil,
				nil,
				dependency.Requirements{
					mockdependency.Requirement("first", "handled as the first", nil),
				},
				nil,
				dependency.ErrNotHandled,
			),
			"two": mockcomponent.NewComponentWDependencies(
				"two",
				nil,
				nil,
				dependency.Requirements{
					mockdependency.Requirement("second", "handled as the second", nil),
				},
				nil,
				dependency.ErrNotHandled,
			),
			"three": mockcomponent.NewComponentWDependencies(
				"three",
				nil,
				nil,
				dependency.Requirements{},
				mockdependency.Dependency("only", "should get used for all requirements", nil, nil),
				nil,
			),
		},
	}

	if err := cl.Validate(ctx); err != nil {
		t.Errorf("unexpected error occurred in validation: %s \n %+v", err.Error(), cl)
	}

	if cl_d := cl.Debug(ctx); cl_d == nil {
		t.Errorf("project debug was empty: %+v", cl_d)
	}
}
