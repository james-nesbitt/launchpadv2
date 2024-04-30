package action_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
)

func Test_EventAddMerge(t *testing.T) {
	es := action.Events{}

	es.Add(&action.Event{Id: "one"})

	if !es.Contains("one") {
		t.Errorf("events says it doesn't contain an added event: %+v", es)
	}

	esa := action.Events{
		"two":   &action.Event{Id: "two"},
		"three": &action.Event{Id: "three"},
	}

	es.Merge(esa)

	if len(es) != 3 {
		t.Errorf("events had wrong length after merge: %+v", es)
	}
}

func Test_EventContains(t *testing.T) {
	es := action.Events{}

	if es.Contains("aaaa") {
		t.Errorf("empty events says it contains an event: %+v", es)
	}

	es.Add(&action.Event{Id: "add"})

	if !es.Contains("add") {
		t.Errorf("events says it doesn't contain an added event: %+v", es)
	}

	if es.Contains("") {
		t.Errorf("events says it contains an event with no id: %+v", es)
	}
}
