package action

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// DeliversEvents this action delivers events.
type DeliversEvents interface {
	DeliversEvents(context.Context) Events
}

// RequiresEvents indicate events that must occur before and after actions for this object
//
//	"before" means before me
//	"after"  means after me
type RequiresEvents interface {
	RequiresEvents(context.Context) (before Events, after Events)
}

// DependencyDeliversEvents add events of a certain key from a map of dependencies to a set of Events.
func DependencyDeliversEvents(ctx context.Context, eventKeys []string, deps map[string]dependency.Dependency, delivers Events) error {
	// collect all events for all delivered dependencies
	errs := []error{}

	for key, d := range deps {
		if d == nil {
			errs = append(errs, fmt.Errorf("%w; %s dependency empty", dependency.ErrNotHandled, key))
			continue
		}

		ddes, ok := d.(DeliversEvents)
		if !ok {
			errs = append(errs, fmt.Errorf("%w; %s dependency does not deliver events", dependency.ErrNotHandled, key))
			continue
		}

		es := ddes.DeliversEvents(ctx)

		for _, key := range eventKeys {
			e := es.Get(key)
			if e == nil {
				continue
			}
			delivers.Add(e)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// DependencyRequiresEvents take events from Requirement Dependency instances and add them to "before" and "after" lists
//
//	There are many cases where you want to be able to determine all of the after and before events to track for a list
//	of dependencies. This helper will build those lists for you from a map of dependency.Requirements
func DependencyRequiresEvents(ctx context.Context, reqs map[string]dependency.Requirement, before, after Events) error {
	errs := []error{}

	for key, r := range reqs {
		if r == nil {
			errs = append(errs, fmt.Errorf("%w; %s requirement not created", dependency.ErrRequirementNotCreated, key))
			continue
		}

		d := r.Matched(ctx)
		if d == nil {
			errs = append(errs, fmt.Errorf("%w; %s requirement not matched", dependency.ErrRequirementNotMatched, key))
			continue
		}

		ddes, ok := d.(DeliversEvents)
		if !ok {
			slog.DebugContext(ctx, fmt.Sprintf("Dependency required doesn't deliver events for %s", key), slog.Any("requirement", r))
			continue
		}

		es := ddes.DeliversEvents(ctx)

		slog.DebugContext(ctx, fmt.Sprintf("Dependency required events found for %s", key), slog.Any("requirement", r), slog.Any("events", es))

		if be := es.Get(dependency.EventKeyActivated); be != nil {
			before.Add(be)
		}
		if ae := es.Get(dependency.EventKeyDeActivated); ae != nil {
			after.Add(ae)
		}

	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// NewEvents from a list of Event instances.
func NewEvents(ne ...*Event) Events {
	es := Events{}
	es.Add(ne...)
	return es
}

type Events map[string]*Event

func (es Events) Add(nes ...*Event) {
	for _, ne := range nes {
		es[ne.Id] = ne
	}
}

func (es Events) Merge(nes Events) {
	for _, ne := range nes {
		es[ne.Id] = ne
	}
}

func (es Events) Contains(id string) bool {
	_, ok := es[id]
	return ok
}

func (es Events) Get(id string) *Event {
	e, _ := es[id]
	return e
}

type Event struct {
	Id string
}
