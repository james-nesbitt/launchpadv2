package righost

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/k0sproject/rig/v2"
)

// NewRigHostHost Host constructor from HostConfig.
func NewHost(hc Config, rig *rig.ClientWithConfig) Host {
	return Host{
		Config: hc,
		state:  &State{},

		rig: rig,
	}
}

// Host definition for one cluster target.
type Host struct {
	Config
	state *State

	rig *rig.ClientWithConfig
}

// Id unique ID for the host.
func (h Host) Id() string {
	return h.Config.Id
}

// HasRole answer boolean if a host has a role.
func (h Host) HasRole(role string) bool {
	for _, r := range h.Config.Roles {
		if r == role {
			return true
		}
	}

	return false
}

// IsWindows.
func (h Host) IsWindows(ctx context.Context) bool {
	return h.rig.IsWindows()
}

// Connect prove that the host is reachable.
func (h Host) Connect(ctx context.Context) error {
	return h.rig.Connect(ctx)
}

// Exec a command.
func (h *Host) Exec(ctx context.Context, cmd string, inr io.Reader, opts host.ExecOptions) (string, string, error) {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Exec: %s", h.Id(), cmd))

	outs := strings.Builder{}
	outw := io.MultiWriter(
		LogCmdLogger{level: "out", host: h},
		&outs,
	)

	errs := strings.Builder{}
	errw := io.MultiWriter(
		LogCmdLogger{level: "err", host: h},
		&errs,
	)

	var rig *rig.Client = h.rig.Client
	if opts.Sudo {
		rig = h.rig.Sudo()
	}

	if cerr := rig.Connect(ctx); cerr != nil {
		return outs.String(), errs.String(), cerr
	}

	w, serr := rig.StartProcess(ctx, cmd, inr, outw, errw)
	if serr != nil {
		return outs.String(), errs.String(), serr
	}

	if werr := w.Wait(); werr != nil {
		return outs.String(), errs.String(), werr
	}

	return outs.String(), errs.String(), nil
}

// InstallPackages install some packages.
func (h *Host) InstallPackages(ctx context.Context, packages []string) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: RIG install packages", h.Id()), slog.Any("packages", packages))
	return h.rig.Sudo().PackageManagerService.PackageManager().Install(ctx, packages...)
}

// Upload from an io.Reader to a file path on the Host.
func (h *Host) Upload(ctx context.Context, src io.Reader, dst string, fm fs.FileMode, opts host.ExecOptions) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Upload: %s", h.Id(), dst))

	rigc := h.rig.Client
	if opts.Sudo {
		rigc = rigc.Sudo()
	}

	fs := rigc.FS()

	bs, rerr := io.ReadAll(src)
	if rerr != nil {
		return rerr
	}

	if err := fs.WriteFile(dst, bs, fm); err != nil {
		return err
	}

	return nil
}

// ServiceEnable enable and start some services.
func (h *Host) ServiceEnable(ctx context.Context, services []string) error {
	rigc := h.rig.Sudo()

	errs := []error{}
	for _, sn := range services {
		s, sgerr := rigc.Service(sn)
		if sgerr != nil {
			errs = append(errs, sgerr)
			continue
		}

		if err := s.Enable(ctx); err != nil {
			errs = append(errs, sgerr)
			continue
		}
		if err := s.Start(ctx); err != nil {
			errs = append(errs, sgerr)
			continue
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// ServiceRestart stop and restart some services.
func (h *Host) ServiceRestart(ctx context.Context, services []string) error {
	rigc := h.rig.Sudo()

	errs := []error{}
	for _, sn := range services {
		s, sgerr := rigc.Service(sn)
		if sgerr != nil {
			errs = append(errs, sgerr)
			continue
		}

		if err := s.Restart(ctx); err != nil {
			errs = append(errs, sgerr)
			continue
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
