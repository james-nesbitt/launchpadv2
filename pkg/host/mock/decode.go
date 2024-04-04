package mock

import (
	"github.com/Mirantis/launchpad/pkg/host"
)

func init() {
	// Register the mock host with the host decoders list
	host.RegisterDecoder("mock", func(d func(interface{}) error) (host.Host, error) {
		var h ho
		err := d(&h)
		return h, err
	})
}
