package righost

import (
	"fmt"
	"log/slog"

	"github.com/k0sproject/rig/v2"

	"github.com/Mirantis/launchpad/pkg/host"
)

func init() {
	host.RegisterDecoder("rig", Decode)
}

func Decode(id string, d func(interface{}) error) (host.Host, error) {
	var c Config
	if err := d(&c); err != nil {
		return nil, err
	}
	c.Id = id

	var rc rig.ClientWithConfig
	if err := d(&rc); err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprintf("decoding rig host: %s", id), slog.Any("config", c), slog.Any("rig", &rc))

	rh := NewHost(c, &rc)
	return &rh, nil
}
