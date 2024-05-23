package project_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/Mirantis/launchpad/pkg/project"
)

func Test_ProjectComponentValidation(t *testing.T) {
	ctx := context.Background()
	e := errors.New("validation error")
	cp := mock.NewComponent(
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
			"one": mock.NewComponentWDependencies(
				"one",
				nil,
				nil,
				dependency.Requirements{
					mock.Requirement("first", "handled as the first", nil),
				},
				nil,
				dependency.ErrNotHandled,
			),
			"two": mock.NewComponentWDependencies(
				"two",
				nil,
				nil,
				dependency.Requirements{
					mock.Requirement("second", "handled as the second", nil),
				},
				nil,
				dependency.ErrNotHandled,
			),
			"three": mock.NewComponentWDependencies(
				"three",
				nil,
				nil,
				dependency.Requirements{},
				mock.Dependency("only", "should get used for all requirements", nil, nil),
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
