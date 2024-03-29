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

// Decoder function for building Product Components, which matches both the yaml and json unmarshallers
type Decoder func(interface{}) error

var (
	// ProductDecoders handlers which can build Product components based on a key type
	ProductDecoders = map[string]func(Decoder) (component.Component, error){
		"k0s":  DecodeK0sComponent,
		"mcr":  DecodeMCRComponent,
		"mke3": DecodeMKE3Component,
		"mke4": DecodeMKE4Component,
		"msr2": DecodeMSR2Component,
		"msr3": DecodeMSR3Component,
		"msr4": DecodeMSR4Component,
	}
)

// DecodeKnownProduct Component of a specific from a decoder such as a yaml decoder
func DecodeKnownProduct(t string, d Decoder) (component.Component, error) {
	dph, ok := ProductDecoders[t]

	if !ok {
		return nil, fmt.Errorf("Product '%s' has no registered product builder", t)
	}

	return dph(d)
}

func DecodeK0sComponent(d Decoder) (component.Component, error) {
	c := k0s.Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'K0s' : %w", err)
	}

	return k0s.NewK0S(c), nil
}

func DecodeMCRComponent(d Decoder) (component.Component, error) {
	c := mcr.Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MCR' : %w", err)
	}

	return mcr.NewMCR(c), nil
}

func DecodeMKE3Component(d Decoder) (component.Component, error) {
	c := mke3.Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MKE3' : %w", err)
	}

	return mke3.NewMKE3(c), nil
}

func DecodeMKE4Component(d Decoder) (component.Component, error) {
	c := mke4.Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MKE4' : %w", err)
	}

	return mke4.NewMKE4(c), nil
}

func DecodeMSR2Component(d Decoder) (component.Component, error) {
	c := msr2.Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MSR2' : %w", err)
	}

	return msr2.NewMSR2(c), nil
}

func DecodeMSR3Component(d Decoder) (component.Component, error) {
	c := msr3.Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MSR3' : %w", err)
	}

	return msr3.NewMSR3(c), nil
}

func DecodeMSR4Component(d Decoder) (component.Component, error) {
	c := msr4.Config{}

	if err := d(&c); err != nil {
		return nil, fmt.Errorf("Failure to unmarshal product 'MSR4' : %w", err)
	}

	return msr4.NewMSR4(c), nil
}
