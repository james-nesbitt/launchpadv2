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

// HostsDepencency a Depencency that provides hosts
type HostsDependency interface {
	ProduceHosts(context.Context) (Hosts, error)
}

// NewHostsDependency constructor for a hostsDep using a factory for retrieving hosts
func NewHostsDependency(id string, describe string, factory func(context.Context) (Hosts, error)) *hostsDep {
	return &hostsDep{
		id:       id,
		describe: describe,
		factory:  factory,
	}
}

type hostsDep struct {
	id       string
	describe string
	factory  func(context.Context) (Hosts, error)

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

// Met is the dependency fulfilled, or is it still blocked/waiting for fulfillment.
func (hd hostsDep) Met(ctx context.Context) error {
	_, err := hd.factory(ctx)
	return err
}

// Validate the the dependency thinks it has what it needs to fulfill it
func (hd hostsDep) Validate(context.Context) error {
	if hd.factory == nil {
		return ErrNoHostsFunction
	}

	return nil
}

// ProduceHosts is the dependency fulfilled, or is it still blocked/waiting for fulfillment.
func (hd hostsDep) ProduceHosts(ctx context.Context) (Hosts, error) {
	return hd.factory(ctx)
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
