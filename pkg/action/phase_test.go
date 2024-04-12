package action_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/dependency/mock"
)

func Test_PhaseOrder(t *testing.T) {
	dA := mock.Dependency("A", "first-dependency", nil, nil)
	dB := mock.Dependency("B", "second-dependency", nil, nil)
	dC := mock.Dependency("C", "third-dependency", nil, nil)
	dD := mock.Dependency("D", "fourth-dependency", nil, nil)

	ps := &action.Phases{
		action.Phase{ // 0 -> 0 :: before "2"
			Id:       "0",
			Provides: dependency.Dependencies{dA}, // required by "2"
			Requires: dependency.Dependencies{},
		},
		action.Phase{ // 1 -> 1 :: before "5", after "0"
			Id:       "1",
			Provides: dependency.Dependencies{dC}, // also provided by "4"
			Requires: dependency.Dependencies{dA}, // provided by "0"
		},
		action.Phase{ // 2 -> 5 :: after "3" and "4"
			Id:       "5",
			Provides: dependency.Dependencies{},
			Requires: dependency.Dependencies{dD, dC}, // provided by "3" and "4"
		},
		action.Phase{ // 3 -> 2 :: before "3" and "4", after "0"
			Id:       "2",
			Provides: dependency.Dependencies{dB}, // required by "3" and "4"
			Requires: dependency.Dependencies{dA}, // provided by "0"
		},
		action.Phase{ // 4 -> 3 :: before "5", after "2"
			Id:       "3",
			Provides: dependency.Dependencies{dD}, // required by "5"
			Requires: dependency.Dependencies{dB}, // provided by "2"
		},
		action.Phase{ // 5 -> 4 :: after "2"
			Id:       "4",
			Provides: dependency.Dependencies{dC}, // required by "5"
			Requires: dependency.Dependencies{dB}, // provided by "2"
		},
	}

	if err := action.OrderPhases(ps, false); err != nil {
		t.Errorf("Phase order error: \n %s", err.Error())
	}

	for id, pos := range map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
	} {
		if (*ps)[pos].Id != id {
			t.Errorf("Phase order incorrect: %s is not in position %d (%s is)", id, pos, (*ps)[pos].Id)
		}
	}
}
