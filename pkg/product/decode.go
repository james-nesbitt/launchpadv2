package product

import (
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/component"
)

var (
	ErrNoProductDecodersRegistered = errors.New("no product decoders have been registered")
	ErrProductDecoderNotRegistered = errors.New("unrecognized product type")
)

var (
	// ProductDecoders handlers which can build Product components based on a key type
	ProductDecoders = map[string]func(func(interface{}) error) (component.Component, error){}
)

// RegisterDecoder register a new host decoder
func RegisterDecoder(k string, d func(func(interface{}) error) (component.Component, error)) {
	ProductDecoders[k] = d
}

// DecodeKnownProduct Component of a specific from a decoder such as a yaml decoder
func DecodeKnownProduct(t string, d func(interface{}) error) (component.Component, error) {
	if len(ProductDecoders) == 0 {
		return nil, ErrNoProductDecodersRegistered
	}

	dph, ok := ProductDecoders[t]

	if !ok {
		return nil, fmt.Errorf("%w; Product '%s' has no registered product builder", ErrProductDecoderNotRegistered, t)
	}

	return dph(d)
}
