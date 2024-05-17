package mock

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

func ProvidesDependencies(dep dependency.Dependency, err error) fillDep {
	return fillDep{
		d:   dep,
		err: err,
	}
}

type fillDep struct {
	d   dependency.Dependency
	err error
}

func (mpd fillDep) ProvidesDependencies(ctx context.Context, r dependency.Requirement) (dependency.Dependency, error) {
	return mpd.d, mpd.err
}

func Dependency(id, description string, validate error, met error) dependency.Dependency {
	return dep{
		id: id,
		d:  description,
		v:  validate,
		m:  met,
	}
}

type dep struct {
	id string
	d  string // describe
	v  error  // validate
	m  error  // met
}

func (md dep) Id() string {
	return md.id
}
func (md dep) Describe() string {
	return md.d
}
func (md dep) Validate(context.Context) error {
	return md.v
}
func (md dep) Met(context.Context) error {
	return md.m
}
