package kubernetes

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

type KubernetesRequirement interface {
	RequiresKubernetes(context.Context) Version
}

func NewKubernetesRequirement(id string, describe string, version Version) *kubeReq {
	return &kubeReq{
		id:      id,
		desc:    describe,
		version: version,
	}
}

// KubernetesRequires dependency requirement for Kubernetes.
type kubeReq struct {
	id   string
	desc string

	version Version

	d KubernetesDependency
}

// Id.
func (kr kubeReq) Id() string {
	return kr.id
}

// Describe.
func (kr kubeReq) Describe() string {
	return kr.desc
}

// Validate.
func (kr kubeReq) Validate(context.Context) error {
	if err := IsValidKubernetesVersion(kr.version); err != nil {
		return err
	}

	return nil
}

// RequiresKubernetes makes the KubernetesRequirement interface.
func (kr kubeReq) RequiresKubernetes(context.Context) Version {
	return kr.version
}

// Match provide a Dependency.
func (kr *kubeReq) Match(d dependency.Dependency) error {
	kd, ok := d.(KubernetesDependency)
	if !ok {
		return fmt.Errorf("Wrong type of dependency given to KubernetesRequirement: %+v", d)
	}

	kr.d = kd
	return nil
}

// Matched confirm that the requirement is Matched.
func (kr kubeReq) Matched(_ context.Context) dependency.Dependency {
	if kr.d == nil {
		return nil
	}

	return kr.d.(dependency.Dependency)
}

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
