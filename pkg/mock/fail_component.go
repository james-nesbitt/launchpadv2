package mock

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/project"
)

type FailComponent struct {
	id   string
	name string
}

func NewFailComponent(id string) *FailComponent {
	return &FailComponent{
		id:   id,
		name: "fail-component",
	}
}

// Name of the Component.
func (c *FailComponent) Name() string {
	return c.name
}

// Debug output.
func (c *FailComponent) Debug() any {
	return c.name
}

// Validate always okay.
func (c *FailComponent) Validate(context.Context) error {
	return nil
}

func (c *FailComponent) CommandBuild(ctx context.Context, cmd *action.Command) error {
	if cmd.Key == project.CommandKeyApply {
		p := &FailPhase{
			id: "fail-phase",
		}
		cmd.Phases.Add(p)
	}
	return nil
}

type FailPhase struct {
	id string
}

func (p *FailPhase) ID() string {
	return p.id
}

func (p *FailPhase) Run(ctx context.Context) error {
	return fmt.Errorf("this phase always fails for testing AI-driven troubleshooting")
}

func init() {
	// Register the fail component in a Way that it can be used for tests
	// But since this is a mock, maybe I'll just use it directly in a test
}
