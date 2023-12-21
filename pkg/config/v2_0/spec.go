package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/cluster"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/product"
	"github.com/Mirantis/launchpad/pkg/product/mcr"
	"github.com/Mirantis/launchpad/pkg/product/mke"
	"github.com/Mirantis/launchpad/pkg/product/msr"
)

// Config full cluster configuration object for the launchpad object
type Config struct {
	Spec *ConfigSpec `yaml:"spec"`
}

// Cluster from this config object
func (c Config) Cluster() (cluster.Cluster, error) {
	return cluster.Cluster{
		Hosts:    c.Spec.Hosts,
		Products: c.Spec.Products.products,
	}, nil
}

// ConfigSpec defines cluster spec.
type ConfigSpec struct {
	Hosts    host.Hosts         `yaml:"hosts" validate:"required,dive,min=1"`
	Products ConfigSpecProducts `yaml:"products,omitempty"`
}

// ConfigProducts a set of products for a cluster
type ConfigSpecProducts struct {
	products map[string]product.Product
}

// UnmarshalYAML custom un-marshaller
func (cp *ConfigSpecProducts) UnmarshalYAML(py *yaml.Node) error {
	cpt := map[string]interface{}{}
	cp.products = map[string]product.Product{}

	if err := py.Decode(&cpt); err != nil {
		return fmt.Errorf("Failed to parse products: %w", err)
	}

	for t, cptp := range cpt {
		if pn, ok := cptp.(yaml.Node); ok {
			switch t {
			case "mcr":
				c := mcr.Config{}

				if err := pn.Decode(&c); err != nil {
					return fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
				}

				cp.products[t] = mcr.NewMCR(c)

			case "mke":
				c := mke.Config{}

				if err := pn.Decode(&c); err != nil {
					return fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
				}

				cp.products[t] = mke.NewMKE(c)

			case "msr":
				c := msr.Config{}

				if err := pn.Decode(&c); err != nil {
					return fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
				}

				cp.products[t] = msr.NewMSR(c)

			default:
				return fmt.Errorf("Product '%s' is not recognized.", t)
			}
		}
	}

	return nil
}
