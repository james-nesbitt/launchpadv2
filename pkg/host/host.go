package host

import (
	"context"
	"io"
	"io/fs"
)

type Host interface {
	// Id return the host unique identifier
	Id() string
	// HasRole answer if the host has the requested role
	HasRole(string) bool

	// IsWindows check, could require a connection
	IsWindows(context.Context) bool

	// Connect to the host to confirm that a connection is possible
	Connect(ctx context.Context) error

	// Exec execute a command
	Exec(ctx context.Context, cmd string, inr io.Reader, opts ExecOptions) (string, string, error)

	// Install the passed system packages
	InstallPackages(ctx context.Context, packages []string) error

	// Upload content from a io.Reader to a path on the Host
	Upload(ctx context.Context, src io.Reader, dst string, fm fs.FileMode, opts ExecOptions) error

	// ServiceEnable activate and start system services
	ServiceEnable(ctx context.Context, services []string) error
	// ServiceRestart stop and restart system services
	ServiceRestart(ctx context.Context, services []string) error
	// ServiceDisable stop and disable system services
	ServiceDisable(ctx context.Context, services []string) error

	// Network describe the host networking
	Network(ctx context.Context) (Network, error)
}

// ExecOptions options to configure the execute.
type ExecOptions struct {
	Sudo bool
}

// Network networking information about the Host
type Network struct {
	PublicAddress    string `yaml:"public_address"`
	PrivateInterface string `yaml:"private_interface"`
	PrivateAddress   string `yaml:"private_address"`
}
