package host

import (
	"errors"
	"fmt"
)

var (
	ErrNoHostDecodersRegistered = errors.New("No host decoders have been registered")
	ErrUnknownHostDecodeType    = errors.New("Unknown host type for decoding")
)

var (
	// HostDecoders function handlers which can convert an interface decoder into a type of Host for the string key. Used in DecodeHost, and allows overrides for injection and testing.
	HostDecoders = map[string]func(string, func(interface{}) error) (Host, error){}
)

// RegisterDecoder register a new host decoder.
func RegisterDecoder(k string, d func(string, func(interface{}) error) (Host, error)) {
	HostDecoders[k] = d
}

// DecodeHost create a Host from the registered decode functions.
func DecodeHost(t string, id string, d func(interface{}) error) (Host, error) {
	if len(HostDecoders) == 0 {
		return nil, ErrNoHostDecodersRegistered
	}

	hd, ok := HostDecoders[t]
	if !ok {
		return nil, fmt.Errorf("%w: %s, only %+v", ErrUnknownHostDecodeType, t, HostDecoders)
	}

	return hd(id, d)
}
