package host

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

type runHooksStep struct {
	id string

	hook string
	ds   dependency.Dependencies
}

func (s runHooksStep) Id() string {
	return fmt.Sprintf("%s:runhooks:%s", s.id, s.hook)
}

func (s runHooksStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running host hooks", slog.String("ID", s.Id()))
	return nil
}
