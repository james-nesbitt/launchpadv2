package mke3

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
	"github.com/Mirantis/launchpad/pkg/phase"
	"github.com/Mirantis/launchpad/pkg/product/mke3/implementation"
)

// NewMKE3 constructor for MKE3 from config.
func NewMKE3(c Config) MKE3 {
	return MKE3{
		config: c,
		state:  State{},
	}
}

// MKE3 product implementation.
type MKE3 struct {
	config Config
	state  State
}

// Name for the component
func (_ MKE3) Name() string {
	return "MKE3"
}

// Debug product debug.
func (_ MKE3) Debug() interface{} {
	return nil
}

// Provides a list of string labels which this Component provides,
func (_ MKE3) Provides() []string {
	return []string{}
}

// Validate that the cluster meets the needs of the Product
func (_ MKE3) Validate(context.Context, host.Hosts, component.Components) ([]string, error) {
	return []string{}, nil
}

// Actions to run for a particular string phase
func (_ MKE3) Actions(context.Context, string) (phase.Actions, error) {
	return phase.Actions{}, nil
}

// ImplementKubernetes
func (_ MKE3) ImplementKubernetes(context.Context) (kubernetes.Kubernetes, error) {
	return kubernetes.NewKubernetes(), nil
}

// ImplementK0sApi
func (_ MKE3) ImplementK0SApi(context.Context) (implementation.API, error) {
	return implementation.NewAPI(), nil
}
