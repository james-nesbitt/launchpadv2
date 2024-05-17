package k0s

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/util/download"
)

type installK0sStep struct {
	baseStep
	id string
	d  *download.QueueDownload
}

func (s installK0sStep) Id() string {
	return fmt.Sprintf("%s:k0s-install", s.id)
}

func (s installK0sStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running k0s install step", slog.String("ID", s.id))

	s.d = download.NewQueueDownload(nil)

	hs, hserr := s.c.GetAllHosts(ctx)
	if hserr != nil {
		return fmt.Errorf("K0s install step could not retrieve host list: %s", hserr.Error())
	}

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		dlberr := s.uploadBinaries(ctx, h)
		if dlberr != nil {
			return dlberr
		}

		return nil
	}); err != nil {
		return fmt.Errorf("Error installing k0s: %s", err.Error())
	}

	return nil
}

func (s installK0sStep) uploadBinaries(ctx context.Context, h *host.Host) error {
	slog.InfoContext(ctx, fmt.Sprintf("%s: downloading binaries", h.Id()))

	he := exec.HostGetExecutor(h)
	hp := exec.HostGetPlatform(h)

	url := s.c.config.DownloadURL(hp.Arch(ctx))
	fs, fn, ferr := s.d.Download(ctx, url)
	if ferr != nil {
		return ferr
	}

	slog.InfoContext(ctx, fmt.Sprintf("%s: uploading k0s binary: %s -> %s", h.Id(), url, fn))
	if err := he.Upload(ctx, fs, fn, 0777, exec.ExecOptions{}); err != nil {
		return err
	}

	return nil
}
