package righost

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"regexp"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/k0sproject/rig/v2"
)

var (
	sbinPath = "PATH=/usr/local/sbin:/usr/sbin:/sbin:$PATH"
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
		LogCmdLogger{level: opts.OutputLevel, host: h},
		&outs,
	)

	errs := strings.Builder{}
	errw := io.MultiWriter(
		LogCmdLogger{level: opts.OutputLevel, host: h},
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

// RemovePackages install some packages.
func (h *Host) RemovePackages(ctx context.Context, packages []string) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: RIG install packages", h.Id()), slog.Any("packages", packages))
	return h.rig.Sudo().PackageManagerService.PackageManager().Remove(ctx, packages...)
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

// ServiceDisable stop and disable some services
func (h *Host) ServiceDisable(ctx context.Context, services []string) error {
	rigc := h.rig.Sudo()

	errs := []error{}
	for _, sn := range services {
		s, sgerr := rigc.Service(sn)
		if sgerr != nil {
			errs = append(errs, fmt.Errorf("%s: service '%s' unknown: %s", h.Id(), sn, sgerr.Error()))
			continue
		}

		if err := s.Stop(ctx); err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: service '%s' did not stop properly", h.Id(), sn), slog.Any("host", h))
		}

		if err := s.Disable(ctx); err != nil {
			errs = append(errs, fmt.Errorf("%s: service '%s' disable failed: %s", h.Id(), sn, err.Error()))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// Network retrieve the Network information for the machine
func (h *Host) Network(ctx context.Context) (host.Network, error) {
	n := host.Network{}

	if pi, err := h.resolvePrivateInterface(ctx); err != nil {
		return n, err
	} else {
		n.PrivateInterface = pi
	}

	if ip, err := h.resolvePublicIP(); err != nil {
		return n, err
	} else {
		n.PublicAddress = ip
	}

	if ip, err := h.resolveInternaIP(ctx, n.PrivateInterface, n.PublicAddress); err != nil {
		return n, err
	} else {
		n.PrivateAddress = ip
	}

	return n, nil
}

func (h Host) resolvePrivateInterface(ctx context.Context) (string, error) {
	cmd := fmt.Sprintf(`%s; (ip route list scope global | grep -P "\b(172|10|192\.168)\.") || (ip route list | grep -m1 default)`, sbinPath)
	o, e, err := h.Exec(ctx, cmd, nil, host.ExecOptions{})
	if err != nil {
		return "", fmt.Errorf("could not detect private interface ;%w : %s", err, e)
	}

	re := regexp.MustCompile(`\bdev (\w+)`)
	match := re.FindSubmatch([]byte(o))
	if len(match) == 0 {
		return "", fmt.Errorf("can't find 'dev' in output")
	}
	return string(match[1]), nil
}

func (h Host) resolvePublicIP() (string, error) {
	return h.rig.ConnectionConfig.SSH.Address, nil
}

func (h Host) resolveInternaIP(ctx context.Context, privateInterface string, publicIP string) (string, error) {
	o, e, err := h.Exec(ctx, fmt.Sprintf("%s ip -o addr show dev %s scope global", sbinPath, privateInterface), nil, host.ExecOptions{})
	if err != nil {
		return "", fmt.Errorf("%s: failed to find private interface: %s :: %s", h.Id(), err.Error(), e)
	}

	lines := strings.Split(o, "\n")
	for _, line := range lines {
		items := strings.Fields(line)
		if len(items) < 4 {
			//log.Debugf("not enough items in ip address line (%s), skipping...", items)
			continue
		}

		idx := strings.Index(items[3], "/")
		if idx == -1 {
			//log.Debugf("no CIDR mask in ip address line (%s), skipping...", items)
			continue
		}
		addr := items[3][:idx]

		if addr != publicIP {
			//log.Infof("%s: using %s as private IP", h, addr)
			if net.ParseIP(addr) != nil {
				return addr, nil
			}
		}
	}
	return publicIP, nil
}
