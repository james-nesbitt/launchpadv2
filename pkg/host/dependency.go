package host

import (
	"context"
	"errors"
)

var (
	ErrNoHostsFunction = errors.New("No hosts dependency option")
)

type HostsDependency interface {
	ProduceHosts(context.Context) *Hosts
}

func NewHostsDependency(id string, describe string, factory func(context.Context) (*Hosts, error)) *hostsDep {
	return &hostsDep{
		id:       id,
		describe: describe,
		factory:  factory,
	}
}

type hostsDep struct {
	id       string
	describe string
	factory  func(context.Context) (*Hosts, error)
}

// Id uniquely identify the Dependency
func (hd hostsDep) Id() string {
	return hd.id
}

// Describe the dependency for logging and auditing
func (hd hostsDep) Describe() string {
	return hd.describe
}

// Validate the the dependency thinks it has what it needs to fullfill it
//
//	Validation does not meet the dependency, it only confirms that the dependency can be met
//	when it is needed.  Each requirement will have its own interface for individual types of
//	requirements.
func (hd hostsDep) Validate(context.Context) error {
	if hd.factory == nil {
		return ErrNoHostsFunction
	}

	return nil
}

// Met is the dependency fullfilled, or is it still blocked/waiting for fulfillment
func (hd hostsDep) Met(ctx context.Context) error {
	_, err := hd.factory(ctx)
	return err
}

// ProduceHosts is the dependency fullfilled, or is it still blocked/waiting for fulfillment
func (hd hostsDep) ProduceHosts(ctx context.Context) *Hosts {
	hs, _ := hd.factory(ctx)
	return hs
}
