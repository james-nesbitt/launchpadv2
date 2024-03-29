package host

import (
	"errors"
	"fmt"

	righost "github.com/Mirantis/launchpad/pkg/host/rig"
)

var (
	ErrUnknownHostDecodeType = errors.New("Unknown host type for decoding")
)

// Decoder function for config which matches both the yaml and json unmarshalling handlers
type Decoder func(interface{}) error

var (
	// HostDecoders function handlers which can convert an interface decoder into a type of Host for the string key. Used in DecodeHost, and allows overrides for injection and testing
	HostDecoders = map[string]func(Decoder) (Host, error){
		"rig": RigHostDecode,
	}
)

// DecodeHost create a Host from the registered decode functions
func DecodeHost(t string, d Decoder) (Host, error) {
	hd, ok := HostDecoders[t]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownHostDecodeType, t)
	}

	return hd(d)
}

func RigHostDecode(d Decoder) (Host, error) {
	var rc righost.Config
	if err := d(&rc); err != nil {
		return nil, err
	}

	rh := righost.NewHost(&rc)
	return &rh, nil
}
