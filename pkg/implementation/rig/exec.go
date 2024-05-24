package rig

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/k0sproject/rig/v2"

	hostexec "github.com/Mirantis/launchpad/pkg/host/exec"
)

// Connect prove that the host is reachable.
func (p *hostPlugin) Connect(ctx context.Context) error {
	return p.rig.Connect(ctx)
}

// Exec a command.
func (p *hostPlugin) Exec(ctx context.Context, cmd string, inr io.Reader, opts hostexec.ExecOptions) (string, string, error) {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Exec: %s", p.hid(), cmd))

	outs := strings.Builder{}
	outw := io.MultiWriter(
		LogCmdLogger{level: opts.OutputLevel, id: p.hid()},
		&outs,
	)

	errs := strings.Builder{}
	errw := io.MultiWriter(
		LogCmdLogger{level: opts.OutputLevel, id: p.hid()},
		&errs,
	)

	var rig *rig.Client = p.rig.Client
	if opts.Sudo {
		rig = p.rig.Sudo()
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
func (p *hostPlugin) InstallPackages(ctx context.Context, packages []string) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: RIG install packages", p.hid()), slog.Any("packages", packages))
	return p.rig.Sudo().PackageManagerService.PackageManager().Install(ctx, packages...)
}

// RemovePackages install some packages.
func (p *hostPlugin) RemovePackages(ctx context.Context, packages []string) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: RIG install packages", p.hid()), slog.Any("packages", packages))
	return p.rig.Sudo().PackageManagerService.PackageManager().Remove(ctx, packages...)
}

// ServiceEnable enable and start some services.
func (p *hostPlugin) ServiceEnable(ctx context.Context, services []string) error {
	rigc := p.rig.Sudo()

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

		if !s.IsRunning(ctx) {
			if err := s.Start(ctx); err != nil {
				errs = append(errs, sgerr)
				continue
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// ServiceRestart stop and restart some services.
func (p *hostPlugin) ServiceRestart(ctx context.Context, services []string) error {
	rigc := p.rig.Sudo()

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

// ServiceDisable stop and disable some services
func (p *hostPlugin) ServiceDisable(ctx context.Context, services []string) error {
	rigc := p.rig.Sudo()

	errs := []error{}
	for _, sn := range services {
		s, sgerr := rigc.Service(sn)
		if sgerr != nil {
			errs = append(errs, fmt.Errorf("%s: service '%s' unknown: %s", p.hid(), sn, sgerr.Error()))
			continue
		}

		if err := s.Stop(ctx); err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: service '%s' did not stop properly", p.hid(), sn), slog.Any("host", p.hid()))
		}

		if err := s.Disable(ctx); err != nil {
			errs = append(errs, fmt.Errorf("%s: service '%s' disable failed: %s", p.hid(), sn, err.Error()))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// ServiceIsRunning is the service running
func (p *hostPlugin) ServiceIsRunning(ctx context.Context, services []string) error {
	rigc := p.rig.Sudo()

	errs := []error{}
	for _, sn := range services {
		s, sgerr := rigc.Service(sn)
		if sgerr != nil {
			errs = append(errs, fmt.Errorf("%s: service '%s' unknown: %s", p.hid(), sn, sgerr.Error()))
			continue
		}

		if !s.IsRunning(ctx) {
			errs = append(errs, fmt.Errorf("%s: service '%s' is not running", p.hid(), sn))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
