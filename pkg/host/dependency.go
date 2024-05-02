package host

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
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

	events dependency.Events
}

// Id uniquely identify the Dependency.
func (hd hostsDep) Id() string {
	return hd.id
}

// Describe the dependency for logging and auditing.
func (hd hostsDep) Describe() string {
	return hd.describe
}

// Validate the the dependency thinks it has what it needs to fulfill it
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

// Met is the dependency fulfilled, or is it still blocked/waiting for fulfillment.
func (hd hostsDep) Met(ctx context.Context) error {
	_, err := hd.factory(ctx)
	return err
}

// ProduceHosts is the dependency fulfilled, or is it still blocked/waiting for fulfillment.
func (hd hostsDep) ProduceHosts(ctx context.Context) *Hosts {
	hs, _ := hd.factory(ctx)
	return hs
}

func (hd *hostsDep) DeliversEvents(ctx context.Context) dependency.Events {
	if hd.events == nil {
		hd.events = dependency.Events{
			dependency.EventKeyActivated: &dependency.Event{
				Id: fmt.Sprintf("%s:%s", hd.Id(), dependency.EventKeyActivated),
			},
			dependency.EventKeyDeActivated: &dependency.Event{
				Id: fmt.Sprintf("%s:%s", hd.Id(), dependency.EventKeyDeActivated),
			},
		}
	}
	return hd.events
}
