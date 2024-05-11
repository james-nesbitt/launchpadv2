package mcr

import (
	"context"
	"fmt"
)

type baseStep struct {
	c *MCR
}

func (s baseStep) Validate(ctx context.Context) error {
	if _, err := s.c.GetAllHosts(ctx); err != nil {
		return fmt.Errorf("MCR Step validation fail: Host validation: %s", err.Error())
	}

	return nil
}
