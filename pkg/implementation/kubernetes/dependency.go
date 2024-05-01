package kubernetes

import (
	"context"
)

// KubernetesDependency Kubernetes providing Dependency.
type KubernetesDependency interface {
	Kubernetes(ctx context.Context) *Kubernetes
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

func (kd *kubeDep) Kubernetes(ctx context.Context) *Kubernetes {
	k, _ := kd.factory(ctx)
	return k
}
