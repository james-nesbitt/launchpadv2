package component

import (
	"errors"
	"fmt"
)

var (
	ErrNoComponentDecodersRegistered = errors.New("no component decoders have been registered")
	ErrComponentDecoderNotRegistered = errors.New("unrecognized component type")
)

var (
	// ComponentDecoders handlers which can build Product components based on a key type.
	ComponentDecoders = map[string]func(string, func(interface{}) error) (Component, error){}
)

// RegisterDecoder register a new host decoder.
func RegisterDecoder(k string, d func(string, func(interface{}) error) (Component, error)) {
	ComponentDecoders[k] = d
}

// DecodeComponent Component of a specific from a decoder such as a yaml decoder.
func DecodeComponent(t string, id string, d func(interface{}) error) (Component, error) {
	if len(ComponentDecoders) == 0 {
		return nil, ErrNoComponentDecodersRegistered
	}

	dph, ok := ComponentDecoders[t]

	if !ok {
		return nil, fmt.Errorf("%w; Component '%s' has no registered component builder", ErrComponentDecoderNotRegistered, t)
	}

	return dph(id, d)
}
