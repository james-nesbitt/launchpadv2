package phase_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/phase"
)

// --- dummy phase actions ---

type dummyAction struct {
	label    string
	validate error
	action   func() error
}

// Label for logging and labelling
func (da *dummyAction) Label() string {
	return da.label
}

// Validate the action
func (da *dummyAction) Validate(context.Context) error {
	return da.validate
}

// Run the Action
func (da *dummyAction) Run(context.Context) error {
	return da.action()
}

// --- tests

func Test_Actions(t *testing.T) {
	var c int
	a := func(_ context.Context) error {
		c = c + 1
		return nil
	}

}
