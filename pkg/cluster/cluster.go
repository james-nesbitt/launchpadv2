package cluster

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/host"
)

var (
	ErrClusterDependenciesNotMet = errors.New("cluster dependencies not met")
)

// Cluster function handler for the complete cluster
type Cluster struct {
	Hosts      host.Hosts
	Components component.Components

	dependencies dependency.Dependencies
}

// Validate the cluster configuration
func (c Cluster) Validate(ctx context.Context) error {

	// Dependency checking
	c.matchDependencies(ctx)

	if uds := c.dependencies.Unmet(); len(uds) > 0 {
		var udses []string
		for _, ud := range uds {
			udses = append(udses, ud.Describe())
		}

		return fmt.Errorf("%w; %s", ErrClusterDependenciesNotMet, strings.Join(udses, "\n"))
	}

	return nil
}

func (c *Cluster) matchDependencies(ctx context.Context) error {
	c.dependencies = dependency.Dependencies{}

	// Collect all the things that provide dependencies
	pds := []dependency.ProvidesDependencies{}

	if cpd, ok := interface{}(c).(dependency.ProvidesDependencies); ok {
		pds = append(pds, cpd)
	}
	for _, h := range c.Hosts {
		if pd, ok := h.(dependency.ProvidesDependencies); ok {
			pds = append(pds, pd)
		}
	}
	for _, c := range c.Components {
		if pd, ok := c.(dependency.ProvidesDependencies); ok {
			pds = append(pds, pd)
		}
	}

	// Run through all source of dependencies, separately, so that we
	// Can do some logging and auditing
	for _, h := range c.Hosts {
		if hd, ok := h.(dependency.HasDependencies); ok {
			for _, r := range hd.Requires(ctx) {
				d := dependency.ProvideRequirementDependency(ctx, r, pds)
				c.dependencies = append(c.dependencies, d)
			}
		}
	}
	for _, cc := range c.Components {
		if hd, ok := cc.(dependency.HasDependencies); ok {
			for _, r := range hd.Requires(ctx) {
				d := dependency.ProvideRequirementDependency(ctx, r, pds)
				c.dependencies = append(c.dependencies, d)
			}
		}
	}

	if ud := c.dependencies.Unmet(); len(ud) > 0 {
		return fmt.Errorf("%w; Unmet depedencies: %+v", ErrClusterDependenciesNotMet, ud)
	}
	return nil
}
