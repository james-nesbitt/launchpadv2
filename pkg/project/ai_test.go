package project

import (
	"context"
	"fmt"
	"testing"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/stretchr/testify/assert"
)

type failPhase struct {
	id string
}

func (p *failPhase) ID() string {
	return p.id
}

func (p *failPhase) Run(ctx context.Context) error {
	return fmt.Errorf("this phase always fails")
}

func TestProjectTroubleshoot(t *testing.T) {
	ctx := context.Background()
	p := &Project{}

	err := p.Troubleshoot(ctx, "project:apply", fmt.Errorf("simulated failure"))
	assert.NoError(t, err)
}

func TestApplyWithAITroubleshoot(t *testing.T) {
	// Normally we would test the CLI RunE function, but it's hard to test
	// We can at least test that Troubleshoot doesn't fail when called
	ctx := context.Background()
	p := &Project{}

	// Setup a project that will fail
	// Using a mock Phase
	cmd := action.NewEmptyCommand("project:apply")
	cmd.Phases.Add(&failPhase{id: "mock-fail-phase"})

	err := cmd.Run(ctx)
	assert.Error(t, err)

	// If the flag were set, we would call troubleshoot
	tserr := p.Troubleshoot(ctx, cmd.Key, err)
	assert.NoError(t, tserr)
}
