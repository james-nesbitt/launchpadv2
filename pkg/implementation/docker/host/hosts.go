package dockerhost

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/Mirantis/launchpad/pkg/host"
)

type HostsOptions struct {
	SudoDocker bool
}

func NewDockerHosts(hosts host.Hosts, opts HostsOptions) Hosts {
	dhs := Hosts{}

	for k, h := range hosts {
		dh := &Host{Host: h}

		if opts.SudoDocker {
			dh.SudoDocker = true
		}

		dhs[k] = dh
	}

	return dhs
}

type Hosts map[string]*Host

// Add Hosts objects to the Hosts set.
func (hs Hosts) Add(ahs ...*Host) {
	for _, ah := range ahs {
		if ah == nil {
			continue
		}
		hs[ah.Id()] = ah
	}
}

// Merge Hosts objects into the Hosts set.
func (hs Hosts) Merge(ahs Hosts) {
	for id, ah := range ahs {
		if ah == nil {
			continue
		}
		hs[id] = ah
	}
}

// Get Host from the Hosts set if it exists, or return nil.
func (hs Hosts) Get(id string) *Host {
	h, _ := hs[id] //nolint: gosimple
	return h
}

// Each Hosts in the set gets the function applied in parallel.
func (hs Hosts) Each(ctx context.Context, f func(context.Context, *Host) error) error {
	errs := []error{}

	wg := sync.WaitGroup{}
	fctx, cancel := context.WithCancel(ctx)

	for _, h := range hs {
		wg.Add(1)
		go func(ctx context.Context, h *Host) {
			err := f(ctx, h)
			if err != nil {
				errs = append(errs, err) // @TODO lock this when writing
			}
			wg.Done()
		}(fctx, h)
	}

	go func() {
		slog.Debug("dockerhosts.Each()")
		wg.Wait()
		cancel()
	}()

	select {
	case <-fctx.Done():
		// Done
		slog.Debug("dockerhosts.Each finished")
	case <-ctx.Done():
		slog.WarnContext(ctx, "dockerhosts.Each() context ended before function calls", slog.Any("hosts", hs))
		cancel()
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
