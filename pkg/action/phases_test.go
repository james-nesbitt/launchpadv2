package action_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	actionmock "github.com/Mirantis/launchpad/pkg/action/mock"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

func Test_PhasesAddMerge(t *testing.T) {
	phs := action.NewPhases()

	var p action.Phase

	p = actionmock.NewPhase("one", nil, nil, nil, nil, nil)

	phs.Add(p)

	if len(phs) != 1 {
		t.Errorf("phases add didn't result in the right number of phases: %+v", phs)
	}

	p = actionmock.NewPhase("two", nil, nil, nil, nil, nil)
	phs.Add(p)

	if len(phs) != 2 {
		t.Errorf("phases 2nd add didn't result in the right number of phases: %+v", phs)
	}

	phs.Merge(action.NewPhases(
		actionmock.NewPhase("three", nil, nil, nil, nil, nil),
		actionmock.NewPhase("four", nil, nil, nil, nil, nil),
	))

	if len(phs) != 4 {
		t.Errorf("phases merge result in the right number of phases: %+v", phs)
	}

	phs.Add(actionmock.NewPhase("four", nil, nil, nil, nil, nil)) // repeats id

	if len(phs) != 4 {
		t.Errorf("phases add with a repeat didn't result in the correct number of phases: %+v", phs)
	}

	for id, p := range phs {
		if id != p.Id() {
			t.Errorf("phases add didn't end up with the right phase keys: %+v", phs)
		}
	}
}

func Test_PhasesContains(t *testing.T) {
	phs := action.NewPhases(
		actionmock.NewPhase("yes", nil, nil, nil, nil, nil),
		actionmock.NewPhase("also", nil, nil, nil, nil, nil),
	)

	if !phs.Contains("yes") {
		t.Errorf("Phases contains failed to report a phase which it should contain: %+v", phs)
	}

	if phs.Contains("no") {
		t.Errorf("Phases contains reports a phase which it shouldn't contain: %+v", phs)
	}

}

func Test_PhasesOrder(t *testing.T) {
	ctx := context.Background()

	// create a bunch of events which we will share across the phases
	ea := dependency.Event{Id: "A"}
	eb := dependency.Event{Id: "B"}
	ec := dependency.Event{Id: "C"}
	ed := dependency.Event{Id: "D"}
	ee := dependency.Event{Id: "E"}
	ef := dependency.Event{Id: "F"}

	// Create a bunch of phases with relationships that create an explicit order:
	//  - no event free phases (as they end up in random order, so not testable)
	//  - every phase has an explicit position related to the rest
	//  - order is intentionally jumbled (change the order to further test)
	//  - NewPhase(Id, <validate func>, <run func>, Delivers, Before, After)
	//     - Delivers means that this phase delivers the event
	//     - Before means that the event must run before this Phase (thereforce the Phase that delivers the event comes before)
	//     - After means that the event must run after this Phase (therefore the Phase that delivers the event comes after)
	ps := action.NewPhases(
		actionmock.NewPhase("6", nil, nil, dependency.NewEvents(&ee, &ef), dependency.NewEvents(&ed), nil),
		actionmock.NewPhase("3", nil, nil, dependency.NewEvents(&eb), dependency.NewEvents(&ea), nil),
		actionmock.NewPhase("5", nil, nil, dependency.NewEvents(&ec, &ed), dependency.NewEvents(&ea, &eb), nil),
		actionmock.NewPhase("7", nil, nil, dependency.NewEvents(&ef), dependency.NewEvents(&ec, &ee), nil),
		actionmock.NewPhase("1", nil, nil, dependency.NewEvents(&ea), nil, dependency.NewEvents(&eb)),
		actionmock.NewPhase("4", nil, nil, nil, dependency.NewEvents(&eb), dependency.NewEvents(&ed, &ee)),
		actionmock.NewPhase("8", nil, nil, nil, dependency.NewEvents(&ef), nil),
		actionmock.NewPhase("0", nil, nil, nil, nil, dependency.NewEvents(&ea, &eb)),
		actionmock.NewPhase("2", nil, nil, dependency.NewEvents(&ea, &eb), nil, nil),
	)

	ops, err := ps.Order(ctx)
	if err != nil {
		t.Fatalf("phase sort fail: %s", err.Error())
	}

	for i, p := range ops {
		if fmt.Sprintf("%d", i) != p.Id() {
			t.Errorf(" phase in wrong order %s != %d : %s", p.Id(), i, phasePrint(ctx, p))
		} else {
			t.Logf("phase in right order %s != %d : %s", p.Id(), i, phasePrint(ctx, p))
		}
	}
}

// helper for pretty printing a Phase relationship to make testing easier.
func phasePrint(ctx context.Context, p action.Phase) string {
	pp := struct {
		Id       string
		Delivers []string
		Before   []string
		After    []string
	}{
		Id:       p.Id(),
		Delivers: []string{},
		Before:   []string{},
		After:    []string{},
	}

	if pd, ok := p.(dependency.DeliversEvents); ok {
		for _, e := range pd.DeliversEvents(ctx) {
			pp.Delivers = append(pp.Delivers, e.Id)
		}
	} else {
		pp.Delivers = append(pp.Delivers, "none")
	}

	if pr, ok := p.(dependency.RequiresEvents); ok {
		bes, aes := pr.RequiresEvents(ctx)
		for _, be := range bes {
			pp.Before = append(pp.Before, be.Id)
		}
		for _, ae := range aes {
			pp.After = append(pp.After, ae.Id)
		}
	} else {
		pp.Before = append(pp.Before, "none")
		pp.After = append(pp.After, "none")
	}

	return fmt.Sprintf("%+v", pp)
}
