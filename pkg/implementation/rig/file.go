package rig

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"

	"github.com/Mirantis/launchpad/pkg/host/exec"
	rig "github.com/k0sproject/rig/v2"
)

// Stat return file info on a file, or an error.
func (p *hostPlugin) Stat(ctx context.Context, path string, opts exec.ExecOptions) (os.FileInfo, error) {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Stat: %s", p.hid(), path))

	var c *rig.Client = p.rig.Client
	if cerr := c.Connect(ctx); cerr != nil {
		return nil, cerr
	}
	if opts.Sudo {
		c = p.rig.Sudo()
		if cerr := c.Connect(ctx); cerr != nil {
			return nil, cerr
		}
	}

	fs := c.FS()

	return fs.Stat(path)
}

// FileExists does a file exist.
func (p *hostPlugin) FileExist(ctx context.Context, path string, opts exec.ExecOptions) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig FileExists: %s", p.hid(), path))

	var c *rig.Client = p.rig.Client
	if cerr := c.Connect(ctx); cerr != nil {
		return cerr
	}
	if opts.Sudo {
		c = p.rig.Sudo()
		if cerr := c.Connect(ctx); cerr != nil {
			return cerr
		}
	}

	fs := c.FS()

	if !fs.FileExist(path) {
		return fmt.Errorf("%w; %s file not found", os.ErrNotExist, path)
	}
	return nil
}

// Upload from an io.Reader to a file path on the Host.
func (p *hostPlugin) Upload(ctx context.Context, src io.Reader, dst string, fm fs.FileMode, opts exec.ExecOptions) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Upload: %s", p.hid(), dst))

	var c *rig.Client = p.rig.Client
	if cerr := c.Connect(ctx); cerr != nil {
		return cerr
	}
	if opts.Sudo {
		c = p.rig.Sudo()
		if cerr := c.Connect(ctx); cerr != nil {
			return cerr
		}
	}

	fs := c.FS()

	bs, rerr := io.ReadAll(src)
	if rerr != nil {
		return rerr
	}

	if err := fs.WriteFile(dst, bs, fm); err != nil {
		return err
	}

	return nil
}

// Download a network addess to a file on the host
//
// @TODO this implementation is currently locked to linux machines with curl installed.
func (p *hostPlugin) Download(ctx context.Context, url string, dst string, fm fs.FileMode, opts exec.ExecOptions) error {
	if _, e, err := p.Exec(ctx, fmt.Sprintf("rm -f %s", dst), nil, opts); err != nil {
		return fmt.Errorf("%s: host couldn't delete the existing K0s binary: %s -> %s :: %s : %s", p.h.Id(), url, dst, err.Error(), e)
	}
	if _, e, err := p.Exec(ctx, fmt.Sprintf("curl --get --silent --location --output %s %s", dst, url), nil, opts); err != nil {
		return fmt.Errorf("%s: host couldn't download the K0s binary: %s -> %s :: %s : %s", p.h.Id(), url, dst, err.Error(), e)
	}
	if _, e, err := p.Exec(ctx, fmt.Sprintf("chmod ug+x %s", dst), nil, opts); err != nil {
		return fmt.Errorf("%s: host couldn't chmod the K0s binary [%s]: %s :: %s", p.h.Id(), dst, err.Error(), e)
	}
	return nil
}

// Rename a file or folder.
func (p *hostPlugin) Rename(ctx context.Context, old, new string, opts exec.ExecOptions) error {
	slog.DebugContext(ctx, fmt.Sprintf("%s: Rig Rename: %s -> %s", p.hid(), old, new))

	var c *rig.Client = p.rig.Client
	if cerr := c.Connect(ctx); cerr != nil {
		return cerr
	}
	if opts.Sudo {
		c = p.rig.Sudo()
		if cerr := c.Connect(ctx); cerr != nil {
			return cerr
		}
	}

	fs := c.FS()

	return fs.Rename(old, new)
}
