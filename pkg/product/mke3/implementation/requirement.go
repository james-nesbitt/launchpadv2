package implementation

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

type MKE3DependencyConfig struct {
	Version string
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

func (mke3r mke3Req) RequiresMKE3(context.Context) MKE3DependencyConfig {
	return mke3r.config
}
