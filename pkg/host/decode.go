package host

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Mirantis/launchpad/pkg/component"
	"gopkg.in/yaml.v3"
)

var (
	ErrUnknownHostPluginDecodeType = errors.New("Unknown host plugin type for decoding")
)

func init() {
	component.RegisterDecoder(ComponentType, DecodeComponent)
}

// DecodeComponent decode a new hosts component from an unmarshall decoder.
func DecodeComponent(id string, d func(interface{}) error) (component.Component, error) {
	chs := NewHosts()
	hc := NewHostsComponent(id, chs)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// decode into a map of hosts, where each host is a map of host plugins
	shsc := map[string]map[string]yaml.Node{}
	if err := d(&shsc); err != nil {
		return hc, fmt.Errorf("failed to decode spec hosts: %s", err.Error())
	}

	if len(shsc) == 0 {
		return hc, fmt.Errorf("no hosts were found when decoding the project")
	}

	errs := []error{}
	for id, shc := range shsc {
		ds := map[string]func(interface{}) error{}
		for k, d := range shc { // collect all of the plugin node .Decode functions
			ds[k] = d.Decode
		}

		h, err := DecodeHost(ctx, id, ds)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", id, err.Error()))
			continue
		}

		chs.Add(h)
	}

	if len(errs) > 0 {
		return hc, errors.Join(errs...)
	}
	return hc, nil
}

// DecodeHost create a Host from the registered iplugin decode functions.
func DecodeHost(ctx context.Context, id string, ds map[string]func(interface{}) error) (*Host, error) {
	h := NewHost(id)

	if len(hostPluginFactories) == 0 {
		return h, ErrNoHostPluginHandlerRegistered
	}

	errs := []error{}

	for t, d := range ds {
		if hperr := h.DecodeHostPlugin(ctx, t, d); hperr != nil {
			errs = append(errs, fmt.Errorf("%s: %s", t, hperr.Error()))
		}
	}

	if len(errs) > 0 {
		return h, fmt.Errorf("%s: error building host plugins: %s", id, errors.Join(errs...).Error())
	}

	return h, nil
}

// DecodeHost create a Host from the registered iplugin decode functions.
func (h *Host) DecodeHostPlugin(ctx context.Context, t string, d func(interface{}) error) error {
	hpf, ok := hostPluginFactories[t]
	if !ok {
		return fmt.Errorf("%w: unknown host plugin type %s :: %+v", ErrUnknownHostPluginDecodeType, t, hostPluginFactories)
	}

	hp, hperr := hpf.HostPluginDecode(ctx, h, d)
	if hperr != nil {
		return fmt.Errorf("%s: error building host plugin: %s", t, hperr.Error())
	}

	slog.Debug(fmt.Sprintf("%s: host plugin: %s", h.Id(), t), slog.Any("plugin", hp))

	h.AddPlugin(hp)

	return nil
}
