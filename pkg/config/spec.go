package config

import (
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/cluster"
)

var (
	ErrNoDecodersRegistered  = errors.New("No spec decoders registered.")
	ErrUnknownSpecDecodeType = errors.New("Unknown config spec type for decoding")
)

var (
	// SpecDecoders function handlers which can convert an interface decoder into a type of Spec for the string key. Used in DecodeSpec, and allows overrides for injection and testing
	SpecDecoders = map[string]func(*cluster.Cluster, func(interface{}) error) error{}
)

// RegisterSpecDecoder
func RegisterSpecDecoder(t string, d func(*cluster.Cluster, func(interface{}) error) error) {
	SpecDecoders[t] = d
}

// DecodeSpec create a Spec from the registered decode functions
func DecodeSpec(t string, cl *cluster.Cluster, d func(interface{}) error) error {
	if len(SpecDecoders) == 0 {
		return ErrNoDecodersRegistered
	}

	hd, ok := SpecDecoders[t]
	if !ok {
		return fmt.Errorf("%w '%s'", ErrUnknownSpecDecodeType, t)
	}

	return hd(cl, d)
}
