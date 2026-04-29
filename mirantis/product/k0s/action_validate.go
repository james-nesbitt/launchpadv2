package k0s

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
)

type validateHostsStep struct {
	baseStep
	id string
}

func (s validateHostsStep) ID() string {
	return fmt.Sprintf("%s:k0s-validate", s.id)
}

// Run validates all hosts satisfy prerequisites before installation begins.
//
// Checks:
//  1. At least one controller host is configured
//  2. All hosts are reachable via exec (connection is established)
//  3. All hostnames are unique (duplicate hostnames cause etcd cluster split-brain)
func (s validateHostsStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s host validation step", slog.String("ID", s.id))

	hs, hserr := s.c.GetAllHosts(ctx)
	if hserr != nil {
		return fmt.Errorf("validation: failed to retrieve hosts: %w", hserr)
	}

	// Check at least one controller.
	controllers, cerr := s.c.GetControllerHosts(ctx)
	if cerr != nil {
		return fmt.Errorf("validation: failed to retrieve controller hosts: %w", cerr)
	}
	if len(controllers) == 0 {
		return fmt.Errorf("validation: no controller hosts configured; at least one controller is required")
	}

	// Check reachability for all hosts in parallel; collect ALL errors.
	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		eh := exec.HostGetExecutor(h)
		if eh == nil {
			// GetAllHosts already filters hosts without an exec plugin, so this
			// path should not be reached in practice.
			return fmt.Errorf("%s: no exec plugin", h.ID())
		}
		if err := eh.Connect(ctx); err != nil {
			return fmt.Errorf("%s: unreachable: %w", h.ID(), err)
		}
		return nil
	}); err != nil {
		// Wrap the joined error so callers see all unreachable hosts at once.
		return fmt.Errorf("host validation failed (connectivity): %w", err)
	}

	// Sequential hostname duplicate check.
	// Runs after connectivity is confirmed so every exec call succeeds.
	hostnames := make(map[string]string) // hostname -> first host ID that reported it
	var dupeErrs []error

	if err := hs.Sequential(ctx, func(ctx context.Context, h *host.Host) error {
		eh := exec.HostGetExecutor(h)
		if eh == nil {
			return nil // already caught above
		}

		stdout, _, err := eh.Exec(ctx, "hostname", nil, exec.ExecOptions{})
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: hostname command failed, skipping duplicate check: %v", h.ID(), err))
			return nil
		}

		hostname := strings.TrimSpace(stdout)
		if hostname == "" {
			slog.WarnContext(ctx, fmt.Sprintf("%s: empty hostname reported", h.ID()))
			return nil
		}

		slog.DebugContext(ctx, fmt.Sprintf("%s: hostname=%s", h.ID(), hostname))

		if first, exists := hostnames[hostname]; exists {
			dupeErrs = append(dupeErrs,
				fmt.Errorf("duplicate hostname %q: hosts %s and %s share the same hostname (causes etcd split-brain)",
					hostname, first, h.ID()))
		} else {
			hostnames[hostname] = h.ID()
		}
		return nil
	}, false); err != nil {
		return fmt.Errorf("validation: hostname check exec error: %w", err)
	}

	if len(dupeErrs) > 0 {
		return fmt.Errorf("host validation failed (duplicate hostnames): %w", errors.Join(dupeErrs...))
	}

	slog.InfoContext(ctx, fmt.Sprintf("host validation passed: %d hosts, %d controllers", len(hs), len(controllers)))
	return nil
}
