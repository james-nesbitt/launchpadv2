package action_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	actionmock "github.com/Mirantis/launchpad/pkg/action/mock"
	"github.com/Mirantis/launchpad/pkg/dependency"
	depmock "github.com/Mirantis/launchpad/pkg/dependency/mock"
)

func Test_CommandBuild(t *testing.T) {
	ctx := context.Background()
	cmd := action.NewEmptyCommand("test")

	chs := []action.CommandHandler{
		actionmock.NewCommandHandler(
			action.NewPhases(
				actionmock.NewPhase("one", nil, nil, nil, nil, nil),
				actionmock.NewPhase("two", nil, nil, nil, nil, nil),
			),
			dependency.Dependencies{
				depmock.Dependency("A", "Dep A", nil, nil),
			},
			dependency.Events{},
			nil,
		),
		actionmock.NewCommandHandler(
			action.NewPhases(
				actionmock.NewPhase("three", nil, nil, nil, nil, nil),
				actionmock.NewPhase("four", nil, nil, nil, nil, nil),
			),
			dependency.Dependencies{
				depmock.Dependency("B", "Dep B", nil, nil),
				depmock.Dependency("C", "Dep C", nil, nil),
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
