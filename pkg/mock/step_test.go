package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action/stepped"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
)

func Test_StepSanity(t *testing.T) {
	id := "test"
	rerr := errors.New("a run error")
	var s stepped.Step = mock.NewStep(id, dependency.Events{}, dependency.Events{}, dependency.Events{}, rerr)

	if s.Id() != id {
		t.Errorf("mock step gave wrong id: %+v", s)
	}
	if err := s.Run(context.Background()); !errors.Is(err, rerr) {
		t.Errorf("mock step gave wrong run error: %+v", s)
	}
}
