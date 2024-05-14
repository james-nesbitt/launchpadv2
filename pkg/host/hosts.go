package host

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

// Hosts collection of Host objects for a cluster.
type Hosts map[string]*Host

// NewHosts Hosts constructer.
func NewHosts(hs ...*Host) Hosts {
	nhs := Hosts{}
	for _, h := range hs {
		nhs[h.Id()] = h
	}
	return nhs
}

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
		slog.DebugContext(ctx, "host.Hosts.Each() waiting")
		wg.Wait()
		cancel()
	}()

	select {
	case <-fctx.Done():
		// Done
		slog.DebugContext(ctx, "hosts.Hosts.Each() finished")
	case <-ctx.Done():
		slog.WarnContext(ctx, "hosts.Hosts.Each() context ended before function calls", slog.Any("hosts", hs))
		cancel()
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
