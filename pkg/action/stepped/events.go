package stepped

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// DeliversEvents indicate what events this Phase delivers.
func (sp SteppedPhase) DeliversEvents(ctx context.Context) dependency.Events {
	es, _ := sp.deliveredEventsCheck(ctx)
	return es
}

func (sp SteppedPhase) deliveredEventsCheck(ctx context.Context) (dependency.Events, error) {
	return dependency.DependencyDeliversEvents(ctx, sp.deps, sp.deks)
}

// RequiresEvents indicate what events this Phase requires (before and after).
func (sp SteppedPhase) RequiresEvents(ctx context.Context) (dependency.Events, dependency.Events) {
	bes, aes, _ := sp.requiredEventsCheck(ctx)
	return bes, aes
}

func (sp SteppedPhase) requiredEventsCheck(ctx context.Context) (dependency.Events, dependency.Events, error) {
	return dependency.DependencyRequiresEvents(ctx, sp.reqs)
}
