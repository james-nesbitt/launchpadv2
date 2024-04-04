package mock_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/dependency"
	mockdependency "github.com/Mirantis/launchpad/pkg/dependency/mock"
)

func Test_MockSanity(t *testing.T) {
	ctx := context.Background()
	var r dependency.Requirement = mockdependency.Requirement(
		"first",
		"handled as the first",
		nil,
	)

	if r.Describe() != "handled as the first" {
		t.Errorf("Unexpected requirement description: %s", r.Describe())
	}

	var d dependency.Dependency = mockdependency.Dependency(
		"handler-first",
		"handle the first dependency",
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
}
