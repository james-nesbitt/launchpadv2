package action_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
)

func Test_CommandBuild(t *testing.T) {
	ctx := t.Context()
	cmd := action.NewEmptyCommand("test")

	chs := []action.CommandHandler{
		mock.NewCommandHandler(
			action.NewPhases(
				mock.NewPhase("one", nil, nil, nil, nil, nil),
				mock.NewPhase("two", nil, nil, nil, nil, nil),
			),
			dependency.Dependencies{
				mock.StaticDependency("A", "Dep A", nil, nil),
			},
			dependency.Events{},
			nil,
		),
		mock.NewCommandHandler(
			action.NewPhases(
				mock.NewPhase("three", nil, nil, nil, nil, nil),
				mock.NewPhase("four", nil, nil, nil, nil, nil),
			),
			dependency.Dependencies{
				mock.StaticDependency("B", "Dep B", nil, nil),
				mock.StaticDependency("C", "Dep C", nil, nil),
			},
			dependency.Events{},
			nil,
		),
	}

	if err := action.BuildCommand(ctx, cmd, chs); err != nil {
		t.Errorf("Build command returned unexpected error: %s", err)
	}

	if len(cmd.Dependencies) != 3 {
		t.Errorf("handled command had the wrong number of dependencies: %+v", cmd.Dependencies)
	}

	if len(cmd.Phases) != 4 {
		t.Errorf("handled command had the wrong number of dependencies")
	}
}
