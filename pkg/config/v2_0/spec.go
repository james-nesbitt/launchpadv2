package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/cluster"
	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/product"
)

// Config full cluster configuration object for the launchpad object
type Config struct {
	Spec *ConfigSpec `yaml:"spec"`
}

// Cluster from this config object
func (c Config) Cluster() (cluster.Cluster, error) {
	return cluster.Cluster{
		Hosts:      c.Spec.Hosts,
		Components: c.Spec.Products.products,
	}, nil
}

// ConfigSpec defines cluster spec.
type ConfigSpec struct {
	Hosts    host.Hosts         `yaml:"hosts" validate:"required,dive,min=1"`
	Products ConfigSpecProducts `yaml:"products,omitempty"`
}

// ConfigHosts a set of hosts for a cluster
type ConfigSpecHosts struct {
	hosts []host.Host
}

// UnmarshalYAML custom un-marshaller
func (ch *ConfigSpecHosts) UnmarshalYAML(py *yaml.Node) error {
	cht := map[string]interface{}{}
	ch.hosts = []host.Host{}

	if err := py.Decode(&cht); err != nil {
		return fmt.Errorf("Failed to parse hosts: %w", err)
	}

	var hpe error // collect all product parsing errors here
	for t, chtp := range cht {
		if hn, ok := chtp.(yaml.Node); ok {
			p, err := host.DecodeHost(t, hn.Decode)
			if err != nil {
				if hpe == nil {
					hpe = fmt.Errorf("Failed to parse hosts %s", "")
				}
				hpe = fmt.Errorf("%w; %s", hpe, err)
				continue
			}

			ch.hosts = append(ch.hosts, p)
		}
	}

	return hpe
}

// ConfigProducts a set of products for a cluster
type ConfigSpecProducts struct {
	products map[string]component.Component
}

// UnmarshalYAML custom un-marshaller
func (cp *ConfigSpecProducts) UnmarshalYAML(py *yaml.Node) error {
	cpt := map[string]interface{}{}
	cp.products = map[string]component.Component{}

	if err := py.Decode(&cpt); err != nil {
		return fmt.Errorf("Failed to parse products: %w", err)
	}

	var ppe error // collect all product parsing errors here
	for t, cptp := range cpt {
		if pn, ok := cptp.(yaml.Node); ok {
			p, err := product.DecodeKnownProduct(t, pn.Decode)
			if err != nil {
				if ppe == nil {
					ppe = fmt.Errorf("Failed to parse products %s", "")
				}
				ppe = fmt.Errorf("%w; %s", ppe, err)
				continue
			}

			cp.products[t] = p
		}
	}

	return ppe
}

type SpecHosts struct {
}
