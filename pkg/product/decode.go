package product

import (
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/product/k0s"
	"github.com/Mirantis/launchpad/pkg/product/mcr"
	"github.com/Mirantis/launchpad/pkg/product/mke3"
	"github.com/Mirantis/launchpad/pkg/product/mke4"
	"github.com/Mirantis/launchpad/pkg/product/msr2"
	"github.com/Mirantis/launchpad/pkg/product/msr3"
	"github.com/Mirantis/launchpad/pkg/product/msr4"
)

// Decoder function for config
type decoder func(interface{}) error

// DecodeKnownProduct Component of a specific from a decoder such as a yaml decoder
func DecodeKnownProduct(t string, d decoder) (component.Component, error) {
	switch t {
	case "k0s":
		c := k0s.Config{}

		if err := d(&c); err != nil {
			return nil, fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
		}

		return k0s.NewK0S(c), nil

	case "mcr":
		c := mcr.Config{}

		if err := d(&c); err != nil {
			return nil, fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
		}

		return mcr.NewMCR(c), nil

	case "mke3":
		c := mke3.Config{}

		if err := d(&c); err != nil {
			return nil, fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
		}

		return mke3.NewMKE3(c), nil
	case "mke4":
		c := mke4.Config{}

		if err := d(&c); err != nil {
			return nil, fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
		}

		return mke4.NewMKE4(c), nil

	case "msr2":
		c := msr2.Config{}

		if err := d(&c); err != nil {
			return nil, fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
		}

		return msr2.NewMSR2(c), nil
	case "msr3":
		c := msr3.Config{}

		if err := d(&c); err != nil {
			return nil, fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
		}

		return msr3.NewMSR3(c), nil
	case "msr4":
		c := msr4.Config{}

		if err := d(&c); err != nil {
			return nil, fmt.Errorf("Failure to unmarshal product '%s' : %w", t, err)
		}

		return msr4.NewMSR4(c), nil

	default:
		return nil, fmt.Errorf("Product '%s' is not recognized.", t)
	}

}
