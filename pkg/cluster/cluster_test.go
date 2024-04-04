package cluster_test

import (
	"context"
	"testing"

	mockcomponent "github.com/Mirantis/launchpad/pkg/component/mock"
	"github.com/Mirantis/launchpad/pkg/dependency"
	mockdependency "github.com/Mirantis/launchpad/pkg/dependency/mock"
	mockhost "github.com/Mirantis/launchpad/pkg/host/mock"

	"github.com/Mirantis/launchpad/pkg/cluster"
	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/phase"
)

func Test_Cluster(t *testing.T) {
	ctx := context.Background()
	cl := cluster.Cluster{
		Hosts: host.Hosts{
			"dummy": mockhost.NewHost(
				[]string{"dummy"},
				nil,
			),
		},
		Components: component.Components{
			"one": mockcomponent.NewComponentWDependencies(
				"one",
				phase.Actions{},
				dependency.Requirements{
					mockdependency.Requirement(
						"first",
						"handled as the first",
						nil,
					),
				},
				nil,
				dependency.ErrDependencyNotHandled,
			),
			"two": mockcomponent.NewComponentWDependencies(
				"two",
				phase.Actions{},
				dependency.Requirements{
					mockdependency.Requirement(
						"second",
						"handled as the second",
						nil,
					),
				},
				nil,
				dependency.ErrDependencyNotHandled,
			),
			"three": mockcomponent.NewComponentWDependencies(
				"three",
				phase.Actions{},
				dependency.Requirements{},
				mockdependency.Dependency(
					"only",
					"should get used for all requirements",
					nil,
				),
				nil,
			),
		},
	}

	if err := cl.Validate(ctx); err != nil {
		t.Errorf("unexpected error occured in validation: %s \n %+v", err.Error(), cl)
	}
}
