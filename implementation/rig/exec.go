package rig

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	rig "github.com/k0sproject/rig/v2"

	"github.com/Mirantis/launchpad/pkg/host/exec"
	hostexec "github.com/Mirantis/launchpad/pkg/host/exec"
)

var (
	ErrConnectFailed = errors.New("failed to connect to rig host")
)

// Connect prove that the host is reachable.
func (p *hostPlugin) Connect(ctx context.Context) error {
	slog.DebugContext(ctx, fmt.Sprintf("rig: connecting to host '%s'", p.h.Id()), slog.Any("rig", p.rig), slog.Any("host", p.h))
	if err := p.rig.Connect(ctx); err != nil {
		return fmt.Errorf("%w; %s", ErrConnectFailed, err.Error())
	}
	return nil
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

	var c *rig.Client = p.rig.Client
	if cerr := c.Connect(ctx); cerr != nil {
		return outs.String(), errs.String(), cerr
	}
	if opts.Sudo {
		c = p.rig.Sudo()
		if cerr := c.Connect(ctx); cerr != nil {
			return outs.String(), errs.String(), cerr
		}
	}

	w, serr := c.StartProcess(ctx, cmd, inr, outw, errw)
	if serr != nil {
		return outs.String(), errs.String(), serr
	}

	wc := make(chan error)
	go func() {
		wc <- w.Wait()
	}()

	select {
	case <-ctx.Done():
		return outs.String(), errs.String(), fmt.Errorf("command cancelled due to context closure")
	case werr := <-wc:
		return outs.String(), errs.String(), werr
	}
}

// ExecInteractive get an interactive shell.
func (p *hostPlugin) ExecInteractive(ctx context.Context, opts exec.ExecOptions) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Exec Interactive", p.hid()))

	var c *rig.Client = p.rig.Client
	if cerr := c.Connect(ctx); cerr != nil {
		return cerr
	}

	return c.ExecInteractive("", os.Stdin, os.Stdout, os.Stderr)
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
		slog.DebugContext(ctx, fmt.Sprintf("enabled service: %s", sn))

		if !s.IsRunning(ctx) {
			if err := s.Start(ctx); err != nil {
				errs = append(errs, sgerr)
				continue
			}
			slog.DebugContext(ctx, fmt.Sprintf("started service: %s", sn))
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

// ServiceDisable stop and disable some services.
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

// ServiceIsRunning is the service running.
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
