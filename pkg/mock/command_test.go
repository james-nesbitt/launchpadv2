package mock_test

import (
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
)

func Test_CommandHandlerSanity(t *testing.T) {
	ctx := t.Context()
	ch1 := mock.NewCommandHandler(
		action.NewPhases(
			mock.NewPhase("one", nil, nil, nil, nil, nil),
			mock.NewPhase("two", nil, nil, nil, nil, nil),
			mock.NewPhase("three", nil, nil, nil, nil, nil),
			mock.NewPhase("four", nil, nil, nil, nil, nil),
		),
		dependency.NewDependencies(
			mock.StaticDependency("A", "Dep A", nil, nil),
			mock.StaticDependency("B", "Dep B", nil, nil),
			mock.StaticDependency("C", "Dep C", nil, nil),
		),
		dependency.Events{},
		nil,
	)

	ch2 := mock.NewCommandHandler(
		action.Phases{},
		dependency.Dependencies{},
		dependency.Events{},
		errors.New("a build error"),
	)

	cmd := action.NewEmptyCommand("my-cmd")

	if err := action.BuildCommand(ctx, cmd, []action.CommandHandler{ch1}); err != nil {
		t.Errorf("error occurred building command using mock handler: %s", err.Error())
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
