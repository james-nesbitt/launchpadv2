package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	actionmock "github.com/Mirantis/launchpad/pkg/action/mock"
	"github.com/Mirantis/launchpad/pkg/dependency"
	depmock "github.com/Mirantis/launchpad/pkg/dependency/mock"
)

func Test_CommandHandlerSanity(t *testing.T) {
	ctx := context.Background()
	ch1 := actionmock.NewCommandHandler(
		action.NewPhases(
			actionmock.NewPhase("one", nil, nil, nil, nil, nil),
			actionmock.NewPhase("two", nil, nil, nil, nil, nil),
			actionmock.NewPhase("three", nil, nil, nil, nil, nil),
			actionmock.NewPhase("four", nil, nil, nil, nil, nil),
		),
		dependency.NewDependencies(
			depmock.Dependency("A", "Dep A", nil, nil),
			depmock.Dependency("B", "Dep B", nil, nil),
			depmock.Dependency("C", "Dep C", nil, nil),
		),
		action.Events{},
		nil,
	)

	ch2 := actionmock.NewCommandHandler(
		action.Phases{},
		dependency.Dependencies{},
		action.Events{},
		errors.New("a build error"),
	)

	cmd := action.NewEmptyCommand("my-cmd")

	if err := action.BuildCommand(ctx, cmd, []action.CommandHandler{ch1}); err != nil {
		t.Errorf("error occured building command using mock handler: %s", err.Error())
	}

	if len(cmd.Phases) != 4 {
		t.Errorf("wrong phase count after build: %+v", cmd.Phases)
	}
	if len(cmd.Dependencies) != 3 {
		t.Errorf("wrong dependency count after build: %+v", cmd.Dependencies)
	}

	if err := action.BuildCommand(ctx, cmd, []action.CommandHandler{ch2}); err == nil {
		t.Error("build did not provided expected error")
	}
}
