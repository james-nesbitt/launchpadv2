package k0s

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/k0sproject/version"
)

type discoverStep struct {
	baseStep
	id string
}

func (s discoverStep) ID() string {
	return fmt.Sprintf("%s:k0s-discover", s.id)
}

// Run discovers per-host k0s state: installed version, running version, role, and cluster ID.
//
// Results are stored in c.state.HostInfo and c.state.Leader.
// Errors from individual hosts are logged as warnings and do not fail the step;
// a host with no k0s binary simply has no InstalledVersion recorded.
func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s discover step", slog.String("ID", s.id))

	hs, hserr := s.c.GetAllHosts(ctx)
	if hserr != nil {
		return fmt.Errorf("discover: failed to retrieve hosts: %w", hserr)
	}

	type result struct {
		hostID string
		info   HostDiscovery
	}

	results := make([]result, 0, len(hs))
	var mu sync.Mutex

	// Discover each host in parallel.
	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		kh := HostGetK0s(h)
		if kh == nil {
			// Host is in the list but has no k0s plugin — already warned by GetAllHosts.
			return nil
		}

		info := HostDiscovery{}

		// Installed binary version.
		if v, verr := kh.Version(ctx); verr == nil {
			pv, perr := version.NewVersion(v.K0s)
			if perr == nil {
				info.InstalledVersion = pv
			} else {
				slog.WarnContext(ctx, fmt.Sprintf("%s: could not parse k0s installed version %q: %s", h.ID(), v.K0s, perr.Error()))
			}
		} else {
			slog.DebugContext(ctx, fmt.Sprintf("%s: k0s binary not found or version unavailable: %s", h.ID(), verr.Error()))
		}

		// Running k0s status.
		if st, sterr := kh.Status(ctx); sterr == nil {
			info.Running = true
			info.Role = st.Role
			info.RunningVersion = st.Version
			slog.InfoContext(ctx, fmt.Sprintf("%s: k0s running as %s version %s", h.ID(), st.Role, st.Version))
		} else {
			slog.DebugContext(ctx, fmt.Sprintf("%s: k0s not running: %s", h.ID(), sterr.Error()))
		}

		mu.Lock()
		results = append(results, result{hostID: h.ID(), info: info})
		mu.Unlock()
		return nil
	}); err != nil {
		// Each() collects errors from all goroutines; log and continue.
		slog.WarnContext(ctx, fmt.Sprintf("discover: some hosts had errors: %s", err.Error()))
	}

	// Populate component state.
	hostInfo := make(map[string]HostDiscovery, len(results))
	for _, r := range results {
		hostInfo[r.hostID] = r.info
	}
	s.c.state.HostInfo = hostInfo

	// Elect leader: pick the first running controller; fall back to first controller.
	controllers, cerr := s.c.GetControllerHosts(ctx)
	if cerr != nil {
		// Controller host retrieval failure means we cannot elect a leader.
		// Propagate the error so the caller knows discovery was incomplete.
		return fmt.Errorf("discover: failed to retrieve controller hosts: %w", cerr)
	}
	if len(controllers) == 0 {
		slog.WarnContext(ctx, "discover: no controller hosts found; cluster may not be functional")
		return nil
	}

	var fallback *host.Host
	for _, h := range controllers {
		if fallback == nil {
			fallback = h
		}
		info, ok := hostInfo[h.ID()]
		if ok && info.Running && info.Role == RoleController {
			s.c.state.Leader = h
			slog.InfoContext(ctx, fmt.Sprintf("discover: elected %s as leader (running controller)", h.ID()))
			return nil
		}
	}

	// No running controller found; fall back to first controller for fresh installs.
	s.c.state.Leader = fallback
	slog.InfoContext(ctx, fmt.Sprintf("discover: no running controller found; using %s as leader for fresh install", fallback.ID()))

	return nil
}
