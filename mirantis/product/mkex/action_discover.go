package mkex

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
)

type discoverStep struct {
	baseStep
	id string
}

func (s discoverStep) Id() string {
	return fmt.Sprintf("%s:%s-discover", s.id, ComponentType)
}

func (s discoverStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running %s discover step", ComponentType, slog.String("ID", s.Id()))

	hs, hserr := s.c.GetAllHosts(ctx)
	if hserr != nil {
		return hserr
	}

	if err := hs.Each(ctx, s.hostOSTreeDiscover); err != nil {
		return err
	}

	return nil
}

func (s discoverStep) hostOSTreeDiscover(ctx context.Context, h *host.Host) error {

	return nil
}
