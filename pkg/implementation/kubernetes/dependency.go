package kubernetes

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

// KubernetesDependency Kubernetes providing Dependency.
type KubernetesDependency interface {
	Kubernetes(ctx context.Context) (*Kubernetes, error)
}

func NewKubernetesDependency(id string, desc string, factory func(context.Context) (*Kubernetes, error)) *kubeDep {
	return &kubeDep{
		id:      id,
		desc:    desc,
		factory: factory,
	}
}

// KubernetesDependency a dependency which can provide a kubernetes client.
type kubeDep struct {
	id   string
	desc string

	factory func(context.Context) (*Kubernetes, error)

	events dependency.Events
}

// Id.
func (kd kubeDep) Id() string {
	return kd.id
}

// Describe.
func (kd kubeDep) Describe() string {
	return kd.desc
}

// Describe.
func (kd kubeDep) Validate(context.Context) error {
	return nil
}

// Describe.
func (kd kubeDep) Met(ctx context.Context) error {
	_, err := kd.factory(ctx)
	return err
}

// DeliversEvents the events that this Dependency can deliver.
func (kd *kubeDep) DeliversEvents(context.Context) dependency.Events {
	if kd.events == nil {
		kd.events = dependency.Events{
			dependency.EventKeyActivated: &dependency.Event{
				Id: fmt.Sprintf("%s:%s", kd.Id(), dependency.EventKeyActivated),
			},
			dependency.EventKeyDeActivated: &dependency.Event{
				Id: fmt.Sprintf("%s:%s", kd.Id(), dependency.EventKeyDeActivated),
			},
		}
	}

	return kd.events
}

// Kubernetes API provide.
func (kd *kubeDep) Kubernetes(ctx context.Context) (*Kubernetes, error) {
	return kd.factory(ctx)
}
