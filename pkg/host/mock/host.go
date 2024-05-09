package mock

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/Mirantis/launchpad/pkg/host"
)

// ---  host ---

func NewHost(id string, roles []string, decodeerror error) ho {
	return ho{
		id:          id,
		Roles:       roles,
		DecodeError: decodeerror,
	}
}

// Host a host type that does nothing, but can be used as a substitute when needed.
type ho struct {
	id string

	Roles []string     `yaml:"roles"`
	Net   host.Network `yaml:"network"`

	DecodeError error
}

func (h ho) Id() string {
	return h.id
}

func (h ho) HasRole(rs string) bool {
	for _, r := range h.Roles {
		if r == rs {
			return true
		}
	}
	return false
}

func (h ho) Connect(ctx context.Context) error {
	return nil
}

func (h ho) Exec(ctx context.Context, cmd string, inr io.Reader, opts host.ExecOptions) (string, string, error) {
	return "see err", "Mock Host was asked to exec", fmt.Errorf("Mock Host don't exec: %+v", h)
}

func (h ho) InstallPackages(ctx context.Context, packages []string) error {
	return fmt.Errorf("Mock Host don't install packages: %+v", h)
}

func (h ho) Upload(ctx context.Context, src io.Reader, dst string, fm fs.FileMode, opts host.ExecOptions) error {
	return fmt.Errorf("Mock Host don't upload: %+v", h)
}

func (h ho) ServiceEnable(ctx context.Context, services []string) error {
	return fmt.Errorf("Mock Host don't manage services: %+v", h)
}

func (h ho) ServiceRestart(ctx context.Context, services []string) error {
	return fmt.Errorf("Mock Host don't manage services: %+v", h)
}

func (h ho) ServiceDisable(ctx context.Context, services []string) error {
	return fmt.Errorf("Mock Host don't manage services: %+v", h)
}

func (h ho) Network(_ context.Context) (host.Network, error) {
	return h.Net, nil
}

func (h ho) IsWindows(_ context.Context) bool {
	return false
}
