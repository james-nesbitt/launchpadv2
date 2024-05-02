package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/mock"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

func Test_PhaseSanity(t *testing.T) {
	ctx := context.Background()
	id := "test"

	var runErr error = errors.New("runError")
	var valErr error = nil

	mp := mock.NewPhase(
		id,
		func(_ context.Context) error { return runErr },
		func(_ context.Context) error { return valErr },
		dependency.NewEvents(&dependency.Event{Id: "A"}),
		dependency.NewEvents(&dependency.Event{Id: "B"}, &dependency.Event{Id: "C"}),
		dependency.NewEvents(),
	)

	var p action.Phase = mp

	if p.Id() != id {
		t.Errorf("Phase returned wrong id: %+v", p)
	}

	if pde, ok := p.(dependency.DeliversEvents); !ok {
		t.Error("mock phase doesn't match DeliverEvents")
	} else if des := pde.DeliversEvents(ctx); len(des) == 0 {
		t.Errorf("mock phase doesn't contain expected Deliver event")
	}

	if pre, ok := p.(dependency.RequiresEvents); !ok {
		t.Error("mock phase doesn't match RequiresEvents")
	} else if bes, aes := pre.RequiresEvents(ctx); len(aes) > 0 {
		t.Errorf("mock phase gave after events when it shouldn't have: %+v", aes)
	} else if len(bes) != 2 {
		t.Errorf("mock phase gave wrong before events")
	}

	if pb, ok := p.(action.Validator); !ok {
		t.Error("mock phase doesn't match Validator")
	} else if err := pb.Validate(ctx); err != valErr {
		t.Error("Phase validate returned unexpected error")
	}

	if err := p.Run(ctx); err != runErr {
		t.Error("Phase run returned unexpected error")
	}
}
