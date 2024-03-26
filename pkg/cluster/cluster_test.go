package cluster_test

import (
	"context"
	"testing"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/phase"
)

func Test_ClusterPhase(t *testing.T) {

}

// DummyComponent
type DummyComponent struct {
	name string
}

// Name for the component
func (dc DummyComponent) Name() string {
	return dc.name
}

// Debug product debug.
func (dc DummyComponent) Debug() interface{} {
	return nil
}

// Provides a list of string labels which this Component provides,
func (dc DummyComponent) Provides() []string {
	return []string{}
}

// Validate that the cluster meets the needs of the Product
func (dc DummyComponent) Validate(context.Context, host.Hosts, component.Components) ([]string, error) {
	return []string{}, nil
}

// Actions to run for a particular string phase
func (dc DummyComponent) Actions(context.Context, string) (phase.Actions, error) {
	return phase.Actions{}, nil
}

// Dummy Implementation

type ImplementsDummy interface {
	ImplementDummy(context.Context) (DummyImplementation, error)
}

// ImplementDummy
func (dc DummyComponent) ImplementDummy(context.Context) (DummyImplementation, error) {
	return DummyImplementation{}, nil
}

// DummyImplementation
type DummyImplementation struct {
}
