package rig

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"

	"github.com/Mirantis/launchpad/pkg/host/exec"
)

// Stat return file info on a file, or an error
func (p *rigHostPlugin) Stat(ctx context.Context, path string, opts exec.ExecOptions) (os.FileInfo, error) {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Stat: %s", p.hid(), path))

	rigc := p.rig.Client
	if opts.Sudo {
		rigc = rigc.Sudo()
	}

	fs := rigc.FS()

	return fs.Stat(path)
}

// FileExists does a file exist
func (p *rigHostPlugin) FileExist(ctx context.Context, path string, opts exec.ExecOptions) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig FileExists: %s", p.hid(), path))

	rigc := p.rig.Client
	if opts.Sudo {
		rigc = rigc.Sudo()
	}

	fs := rigc.FS()

	if !fs.FileExist(path) {
		return fmt.Errorf("%w; %s file not found", os.ErrNotExist, path)
	}
	return nil
}

// Upload from an io.Reader to a file path on the Host.
func (p *rigHostPlugin) Upload(ctx context.Context, src io.Reader, dst string, fm fs.FileMode, opts exec.ExecOptions) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Upload: %s", p.hid(), dst))

	rigc := p.rig.Client
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

// Rename a file or folder
func (p *rigHostPlugin) Rename(ctx context.Context, old, new string, opts exec.ExecOptions) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Rename: %s -> %s", p.hid(), old, new))

	rigc := p.rig.Client
	if opts.Sudo {
		rigc = rigc.Sudo()
	}

	fs := rigc.FS()

	return fs.Rename(old, new)
}
