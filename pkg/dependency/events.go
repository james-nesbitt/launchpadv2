package dependency

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrDependencyEventsFailure = errors.New("could not sort out dependency events")
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
func DependencyDeliversEvents(ctx context.Context, deps Dependencies, eventKeys []string) (Events, error) {
	// collect all events for all delivered dependencies
	errs := []error{}
	des := NewEvents()

	for _, d := range deps {
		if d == nil {
			//errs = append(errs, fmt.Errorf("%w; dependency empty", ErrNotHandled))
			continue
		}

		ddes, ok := d.(DeliversEvents) // Does this dependency deliver events?
		if !ok {
			errs = append(errs, fmt.Errorf("%w; %s dependency does not deliver events", ErrNotHandled, d.Id()))
			continue
		}

		es := ddes.DeliversEvents(ctx)

		if len(eventKeys) == 0 {
			des.Merge(es)
		} else {
			for _, key := range eventKeys {
				e := es.Get(key)
				if e == nil {
					continue
				}
				des.Add(e)
				break
			}
		}
	}

	if len(errs) > 0 {
		return des, errors.Join(errs...)
	}
	return des, nil
}

// DependencyRequiresEvents take events from Requirement Dependency instances and add them to "before" and "after" lists
//
//	There are many cases where you want to be able to determine all of the after and before events to track for a list
//	of dependencies. This helper will build those lists for you from a map of Requirements
func DependencyRequiresEvents(ctx context.Context, reqs Requirements) (Events, Events, error) {
	errs := []error{}
	bes := NewEvents()
	aes := NewEvents()

	for _, r := range reqs {
		if r == nil {
			errs = append(errs, fmt.Errorf("%w; %s requirement is empty", ErrRequirementNotCreated, r.Id()))
			continue
		}

		d := r.Matched(ctx)
		if d == nil {
			errs = append(errs, fmt.Errorf("%w; %s requirement not matched", ErrRequirementNotMatched, r.Id()))
			continue
		}

		ddes, ok := d.(DeliversEvents) // Does this requirement dependnecy deliver any events?
		if !ok {
			slog.DebugContext(ctx, fmt.Sprintf("Dependency required doesn't deliver events for %s", r.Id()), slog.Any("requirement", r))
			continue
		}

		es := ddes.DeliversEvents(ctx)

		slog.DebugContext(ctx, fmt.Sprintf("Dependency required events found for %s", r.Id()), slog.Any("requirement", r), slog.Any("events", es))

		if be := es.Get(EventKeyActivated); be != nil {
			bes.Add(be)
		}
		if ae := es.Get(EventKeyDeActivated); ae != nil {
			aes.Add(ae)
		}

	}

	if len(errs) > 0 {
		return bes, aes, errors.Join(errs...)
	}
	return bes, aes, nil
}

// NewEvents from a list of Event instances.
func NewEvents(ne ...*Event) Events {
	es := Events{}
	es.Add(ne...)
	return es
}

// Events a map of Event objects
//
// @NOTE the event keys are arbitrary if you create the map yourself.
type Events map[string]*Event

func (es Events) Add(nes ...*Event) {
	for _, ne := range nes {
		es[ne.Id] = ne
	}
}

func (es Events) Merge(nes Events) {
	for k, ne := range nes {
		es[k] = ne
	}
}

func (es Events) Contains(id string) bool {
	_, ok := es[id]
	return ok
}

func (es Events) Get(id string) *Event {
	e, _ := es[id] //nolint:gosimple
	return e
}

type Event struct {
	Id string
}
