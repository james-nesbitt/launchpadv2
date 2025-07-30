package project_test

import (
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/Mirantis/launchpad/pkg/project"
)

func Test_ProjectComponentValidation(t *testing.T) {
	ctx := t.Context()
	e := errors.New("validation error")
	cp := mock.NewSimpleComponent(
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
	ctx := t.Context()
	cl := project.Project{
		Components: component.Components{
			"one": mock.NewComponentWDependencies(
				"one",
				nil,
				nil,
				mock.ReqFactoryStatic(dependency.Requirements{
					mock.SimpleRequirement("first", "handled as the first"),
				}),
				mock.DepFactoryStatic(nil, dependency.ErrNotHandled),
			),
			"two": mock.NewComponentWDependencies(
				"two",
				nil,
				nil,
				mock.ReqFactoryStatic(dependency.Requirements{
					mock.SimpleRequirement("second", "handled as the second"),
				}),
				mock.DepFactoryStatic(nil, dependency.ErrNotHandled),
			),
			"three": mock.NewComponentWDependencies(
				"three",
				nil,
				nil,
				mock.ReqFactoryStatic(dependency.Requirements{}),
				mock.DepFactoryStatic(mock.StaticDependency("only", "should get used for all requirements", nil, nil), nil),
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
