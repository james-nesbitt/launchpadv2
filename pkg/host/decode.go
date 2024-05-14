package host

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrNoHostPluginDecodersRegistered = errors.New("No host plugin decoders have been registered")
	ErrUnknownHostPluginDecodeType    = errors.New("Unknown host plugin type for decoding")
)

var (
	// HostPluginDecoders function handlers which can convert an interface decoder into a type of Host for the string key. Used in DecodeHost, and allows overrides for injection and testing.
	HostPluginDecoders = map[string]func(context.Context, *Host, func(interface{}) error) (HostPlugin, error){}
)

// RegisterPluginDecoder register a new host decoder.
func RegisterPluginDecoder(k string, d func(context.Context, *Host, func(interface{}) error) (HostPlugin, error)) {
	HostPluginDecoders[k] = d
}

// DecodeHost create a Host from the registered iplugin decode functions.
func DecodeHost(ctx context.Context, id string, ds map[string]func(interface{}) error) (*Host, error) {
	h := NewHost(id)

	if len(HostPluginDecoders) == 0 {
		return h, ErrNoHostPluginDecodersRegistered
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
	if len(HostPluginDecoders) == 0 {
		return ErrNoHostPluginDecodersRegistered
	}

	hpd, ok := HostPluginDecoders[t]
	if !ok {
		return fmt.Errorf("%w: unknown host plugin type %s :: %+v", ErrUnknownHostPluginDecodeType, t, HostPluginDecoders)
	}

	hp, hperr := hpd(ctx, h, d)
	if hperr != nil {
		return fmt.Errorf("%s: error building host plugin: %s", t, hperr.Error())
	}

	slog.Debug(fmt.Sprintf("%s: host plugin: %s", h.Id(), t), slog.Any("plugin", hp))

	h.AddPlugin(hp)

	return nil
}
