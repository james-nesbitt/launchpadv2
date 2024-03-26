package host

import (
	"errors"
	"fmt"

	righost "github.com/Mirantis/launchpad/pkg/host/rig"
)

var (
	ErrUnknownHostDecodeType = errors.New("Unknown host type for decoding")
)

// Decoder function for config
type decoder func(interface{}) error

// DecodeHost create a Host from a decode function
func DecodeHost(t string, d decoder) (Host, error) {
	switch t {
	case "rig":
		var rc righost.Config
		if err := d(&rc); err != nil {
			return nil, err
		}
		rh := righost.NewHost(&rc)
		return &rh, nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownHostDecodeType, t)
	}
}
