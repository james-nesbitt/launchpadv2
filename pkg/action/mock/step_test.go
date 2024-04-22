package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/action/mock"
)

func Test_StepSanity(t *testing.T) {
	id := "test"
	rerr := errors.New("a run error")
	var s action.Step = mock.NewStep(id, action.Events{}, action.Events{}, action.Events{}, rerr)

	if s.Id() != id {
		t.Errorf("mock step gave wrong id: %+v", s)
	}
	if err := s.Run(context.Background()); !errors.Is(err, rerr) {
		t.Errorf("mock step gave wrong run error: %+v", s)
	}
}
