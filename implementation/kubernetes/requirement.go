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

	d dependency.Dependency
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
	_, ok := d.(KubernetesDependency)
	if !ok {
		return fmt.Errorf("Wrong type of dependency given to KubernetesRequirement: %+v", d)
	}

	kr.d = d
	return nil
}

// Matched confirm that the requirement is Matched.
func (kr kubeReq) Matched(_ context.Context) dependency.Dependency {
	return kr.d
}
