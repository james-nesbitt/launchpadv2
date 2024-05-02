package implementation

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ImplementationType = "mke3-client"
)

// ProvidesMKE3 defines an object that can provides an MKE API instance
//
// @NOTE Any dependency that implements this interface is acceptable for the
//
//	related RequiresMKE3 Requirement
type ProvidesMKE3 interface {
	ProvidesMKE3(ctx context.Context) *API
}

// NewMKE3Dependency create a new Dependency for the MKE3 API.
func NewMKE3Dependency(id string, description string, factory func(context.Context) (*API, error)) *mke3Dep {
	return &mke3Dep{
		id:      id,
		desc:    description,
		factory: factory,
	}
}

type mke3Dep struct {
	id   string
	desc string

	factory func(context.Context) (*API, error)

	events dependency.Events
}

// Id uniquely identify this Dependency.
func (mke3d mke3Dep) Id() string {
	return mke3d.id
}

// Describe this Dependency for logging/auditing.
func (mke3d mke3Dep) Describe() string {
	return mke3d.desc
}

// Validate that this Dependency has what it needs to fulfill.
func (mke3d mke3Dep) Validate(context.Context) error {
	if mke3d.factory == nil {
		return dependency.ErrShouldHaveHandled
	}

	return nil
}

// Met check if this Dependency is Met.
func (mke3d mke3Dep) Met(ctx context.Context) error {
	_, err := mke3d.factory(ctx)
	return err
}

// DeliversEvents indicate that this Dependency will deliver events.
//
// @NOTE the keys for these events indicate the nature of the event.
func (mke3d *mke3Dep) DeliversEvents(context.Context) dependency.Events {
	if mke3d.events == nil {
		mke3d.events = dependency.Events{
			dependency.EventKeyActivated: &dependency.Event{
				Id: fmt.Sprintf("%s:%s", mke3d.Id(), dependency.EventKeyActivated),
			},
			dependency.EventKeyDeActivated: &dependency.Event{
				Id: fmt.Sprintf("%s:%s", mke3d.Id(), dependency.EventKeyDeActivated),
			},
		}
	}

	return mke3d.events
}

// ProvidesMKE3 indicate that this Dependency provides an MKE3 API
//
// @NOTE this meets the ProvidesMKE3 interface.
func (mke3d mke3Dep) ProvidesMKE3(ctx context.Context) *API {
	d, _ := mke3d.factory(ctx)
	return d
}
