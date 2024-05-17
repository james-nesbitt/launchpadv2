package mock_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/mock"
)

func Test_MockDependencySanity(t *testing.T) {
	ctx := context.Background()
	var r dependency.Requirement = mock.Requirement(
		"first",
		"handled as the first",
		nil,
	)

	if r.Describe() != "handled as the first" {
		t.Errorf("Unexpected requirement description: %s", r.Describe())
	}

	var d dependency.Dependency = mock.Dependency(
		"handler-first",
		"handle the first dependency",
		nil,
		nil,
	)

	if d.Describe() != "handle the first dependency" {
		t.Errorf("Unexpected dependency description: %s", d.Describe())
	}

	if err := r.Match(d); err != nil {
		t.Errorf("Requirement match resulted in an error: %+v", r)
	}
	if rd := r.Matched(ctx); rd == nil {
		t.Errorf("Requirement that ws just matched thinks that it isn't matched: %+v", r)
	} else if d.Id() != rd.Id() {
		t.Error("Matched requirement returned the wrong dependency")
	}

	if err := d.Met(ctx); err != nil {
		t.Errorf("Dependency returned unexpected error: %s", err.Error())
	}
}
