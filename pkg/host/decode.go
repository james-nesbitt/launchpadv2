package host

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrUnknownHostPluginDecodeType = errors.New("Unknown host plugin type for decoding")
)

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
