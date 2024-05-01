package cluster_test

import (
	"context"
	"errors"
	"testing"

	mockcomponent "github.com/Mirantis/launchpad/pkg/component/mock"
	"github.com/Mirantis/launchpad/pkg/dependency"
	mockdependency "github.com/Mirantis/launchpad/pkg/dependency/mock"

	"github.com/Mirantis/launchpad/pkg/cluster"
	"github.com/Mirantis/launchpad/pkg/component"
)

func Test_ClusterComponentValidation(t *testing.T) {
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

	cl := cluster.Cluster{
		Components: component.Components{
			"mock": cp,
		},
	}

	if err := cl.Validate(ctx); err == nil {
		t.Errorf("expected error, none returned in validation: %+v", cl)
	}

}

func Test_ClusterDependencyValidation(t *testing.T) {
	ctx := context.Background()
	cl := cluster.Cluster{
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
		t.Errorf("cluster debug was empty: %+v", cl_d)
	}
}
