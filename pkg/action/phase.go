package action

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

type ProvidesPhases interface {
	CommandPhases(ctx context.Context, command string) (dependency.Dependencies, PhaseGenerator)
}

type PhaseGenerator func(ctx context.Context) Phases

type Phases []Phase

type Phase interface {
	Validate(context.Context) error
	Run(context.Context) error
	Rollback(context.Context) error
}

type PhaseHasDependencies interface {
	RequiresDependencies(context.Context) dependency.Dependencies
}

type PhaseFullfillsDependencies interface {
	FullfillsDependencies(context.Context, dependency.Dependency) bool
}
