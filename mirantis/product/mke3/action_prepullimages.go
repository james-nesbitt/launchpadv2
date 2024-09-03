package mke3

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	dockertypesimage "github.com/docker/docker/api/types/image"

	dockerimplementation "github.com/Mirantis/launchpad/implementation/docker"
	dockerhost "github.com/Mirantis/launchpad/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/host"
)

type prePullImagesStep struct {
	baseStep
	id string
}

func (s prePullImagesStep) Id() string {
	return fmt.Sprintf("%s:mke3-image-prepull", s.id)
}

func (s prePullImagesStep) Run(ctx context.Context) error {
	if !s.c.config.imagePrepullRequired() {
		slog.InfoContext(ctx, "MKE3 Image Pre-Pull not required. Skipping.", slog.String("ID", s.Id()))
	}
	slog.InfoContext(ctx, "running MKE3 Image Pre-Pull step", slog.String("ID", s.Id()))

	hs, gerr := s.c.GetManagerHosts(ctx)
	if gerr != nil {
		return fmt.Errorf("MKE3 pre-pull images step could not retrieve manager list: %s", gerr.Error())
	}

	var dpo dockertypesimage.PullOptions = s.c.config.imagePullOptions()

	if err := hs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: prepulling images", h.Id()), slog.Any("host", h))
		is, err := listImages(ctx, h, s.c.config)
		if err != nil {
			return fmt.Errorf("%s error listing images for prepull: %s", h.Id(), err.Error())
		}

		slog.DebugContext(ctx, fmt.Sprintf("%s: images for prepull: %+v", h.Id(), is), slog.Any("host", h), slog.Any("images", is))
		return prePullImages(ctx, h, is, dpo)
	}); err != nil {
		return fmt.Errorf("MKE3 image pre-pull failed: %s", err.Error())
	}

	return nil
}

func listImages(ctx context.Context, h *host.Host, c Config) ([]string, error) {
	is := []string{}

	o, e, err := dockerhost.HostGetDockerExec(h).Run(ctx, []string{"run", c.bootstrapperImage(), "images", "--list"}, dockerimplementation.RunOptions{})
	if err != nil {
		return is, fmt.Errorf("failed to list images: %s : %s", err.Error(), e)
	}

	scanner := bufio.NewScanner(strings.NewReader(o))
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "time=") {
			continue
		}
		is = append(is, l)
	}

	return is, nil
}

func prePullImages(ctx context.Context, h *host.Host, images []string, dpo dockertypesimage.PullOptions) error {
	d := dockerhost.HostGetDockerExec(h)

	errs := []error{}
	for _, i := range images {
		if _, err := d.ImagePull(ctx, i, dpo); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s: image pull failed: %s", h.Id(), errors.Join(errs...).Error())
	}
	return nil
}
