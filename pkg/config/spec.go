package config

import (
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/project"
)

var (
	ErrNoDecodersRegistered  = errors.New("no spec decoders registered")
	ErrUnknownSpecDecodeType = errors.New("unknown config spec type for decoding")
)

// SpecDecoders function handlers which can convert an interface decoder into a type of Spec for the string key. Used in DecodeSpec, and allows overrides for injection and testing.
var SpecDecoders = map[string]func(*project.Project, func(any) error) error{}

// RegisterSpecDecoder register new spec decoveders - to inject new ways of decoding config into launchpad.
func RegisterSpecDecoder(t string, d func(*project.Project, func(any) error) error) {
	SpecDecoders[t] = d
}

// DecodeSpec create a Spec from the registered decode functions.
func DecodeSpec(t string, cl *project.Project, d func(any) error) error {
	if len(SpecDecoders) == 0 {
		return ErrNoDecodersRegistered
	}

	hd, ok := SpecDecoders[t]
	if !ok {
		return fmt.Errorf("%w '%s'", ErrUnknownSpecDecodeType, t)
	}

	return hd(cl, d)
}
