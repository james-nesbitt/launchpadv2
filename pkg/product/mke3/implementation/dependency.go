package implementation

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/action"
	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	ImplementationType = "mke3-client"
)

type MKE3Dependency interface {
	ProvidesMKE3(ctx context.Context) *API
}

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

	events action.Events
}

func (mke3d mke3Dep) Id() string {
	return mke3d.id
}

func (mke3d mke3Dep) Describe() string {
	return mke3d.desc
}

func (mke3d mke3Dep) Validate(context.Context) error {
	if mke3d.factory == nil {
		return dependency.ErrShouldHaveHandled
	}

	return nil
}

func (mke3d mke3Dep) Met(ctx context.Context) error {
	_, err := mke3d.factory(ctx)
	return err
}

func (mke3d *mke3Dep) DeliversEvents(context.Context) action.Events {
	if mke3d.events == nil {
		mke3d.events = action.Events{
			dependency.EventKeyActivated: &action.Event{
				Id: fmt.Sprintf("%s:%s", mke3d.Id(), dependency.EventKeyActivated),
			},
			dependency.EventKeyDeActivated: &action.Event{
				Id: fmt.Sprintf("%s:%s", mke3d.Id(), dependency.EventKeyDeActivated),
			},
		}
	}

	return mke3d.events
}

func (mke3d mke3Dep) ProvidesMKE3(ctx context.Context) *API {
	d, _ := mke3d.factory(ctx)
	return d
}
