package mke3

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/product/mke3/implementation"
)

type MKE3DependencyConfig struct {
	Version string
}

// ValidateDependencyConfig validate that the dependency configuration can be met
func (p MKE3) ValidateMKE3DependencyConfig(dc MKE3DependencyConfig) error {
	return nil
}

// Requires declare that we need a HostsRoles dependency
func (p *MKE3) Requires(_ context.Context) dependency.Requirements {
	p.dhr = dockerhost.NewDockerHostsRequirement(
		p.Name(),
		fmt.Sprintf("%s: Needs docker hosts to install on", ComponentType),
		docker.Version{},
	)

	return dependency.Requirements{
		p.dhr,
	}
}

func (p *MKE3) getDockerHostsDependency(ctx context.Context) (*dockerhost.DockerHosts, error) {
	d := p.dhr.Matched(ctx)
	if d == nil {
		return nil, fmt.Errorf("%s: wanted to get DockerHosts dependency, but it isn't matched", ComponentType)
	}

	if err := d.Met(ctx); err != nil {
		return nil, fmt.Errorf("%s; dockerhosts dependency not Met yet; %w", ComponentType, err)
	}

	dhd, ok := d.(dockerhost.DockerHostsDependency)
	if !ok {
		return nil, fmt.Errorf("%s: dockerhosts dependency is wrong type: %+v", ComponentType, d)
	}

	dh := dhd.ProvidesDockerHost(ctx)

	if dh == nil {
		return nil, fmt.Errorf("%s: dockerhosts dependency gave empty docker hosts: %+v", ComponentType, dhd)
	}

	return dh, nil
}

// Provides dependencies
func (p *MKE3) Provides(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	if mke3r, ok := r.(MKE3Requirement); ok {
		// MKE3 dependency

		c := mke3r.RequiresMKE3(ctx)

		if err := p.ValidateMKE3DependencyConfig(c); err != nil {
			return nil, err
		}

		d := NewMKE3Dependency(
			fmt.Sprintf("%s:%s", ComponentType, implementation.ImplementationType),
			fmt.Sprintf("%s: MKE implementation for : %s", ComponentType, r.Describe()),
			func(_ context.Context) (*MKE3, error) {
				return p, nil
			},
		)

		return d, nil
	}
	if k8sr, ok := r.(kubernetes.KubernetesRequirement); ok {
		// Kubernetes dependency

		c := k8sr.RequiresKubernetes(ctx)

		if err := p.ValidateK8sDependencyConfig(c); err != nil {
			return nil, err
		}

		d := kubernetes.NewKubernetesDependency(
			fmt.Sprintf("%s:%s", ComponentType, kubernetes.ImplementationType),
			fmt.Sprintf("%s: kubernetes implementation: %s", ComponentType, r.Describe()),
			func(ctx context.Context) (*kubernetes.Kubernetes, error) {
				return p.kubernetesImplementation(ctx)
			},
		)

		return d, nil
	}

	return nil, nil
}

type MKE3Requirement interface {
	RequiresMKE3(context.Context) MKE3DependencyConfig
}

func NewMKE3Requirement(id string, desc string, config MKE3DependencyConfig) *mke3Req {
	return &mke3Req{
		id:     id,
		desc:   desc,
		config: config,
	}
}

type MKE3Dependency interface {
	ProvidesMKE3(ctx context.Context) *MKE3
}

func NewMKE3Dependency(id string, description string, factory func(context.Context) (*MKE3, error)) *mke3Dep {
	return &mke3Dep{
		id:      id,
		desc:    description,
		factory: factory,
	}
}

type mke3Req struct {
	id   string
	desc string

	config MKE3DependencyConfig

	dep dependency.Dependency
}

func (mke3r mke3Req) Id() string {
	return mke3r.id
}

func (mke3r mke3Req) Describe() string {
	return mke3r.desc
}

func (mke3r *mke3Req) Match(d dependency.Dependency) error {
	if _, ok := d.(MKE3Dependency); !ok {
		return fmt.Errorf("%w; MKE3 Requirement did not receive an MKE3 Dependency : %+v", dependency.ErrDependencyNotMatched, d)
	}

	mke3r.dep = d
	return nil
}

func (mke3r mke3Req) Matched(context.Context) dependency.Dependency {
	return mke3r.dep
}

func (mke3r mke3Req) RequiresMKE3(_ context.Context) MKE3DependencyConfig {
	return mke3r.config
}

type mke3Dep struct {
	id   string
	desc string

	factory func(context.Context) (*MKE3, error)
}

func (mke3d mke3Dep) Id() string {
	return mke3d.id
}

func (mke3d mke3Dep) Describe() string {
	return mke3d.desc
}

func (mke3d mke3Dep) Validate(context.Context) error {
	if mke3d.factory == nil {
		return dependency.ErrDependencyShouldHaveHandled
	}

	return nil
}

func (mke3d mke3Dep) Met(ctx context.Context) error {
	_, err := mke3d.factory(ctx)
	return err
}

func (mke3d mke3Dep) ProvidesMKE3(ctx context.Context) *MKE3 {
	d, _ := mke3d.factory(ctx)
	return d
}
