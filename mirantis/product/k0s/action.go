package k0s

import (
	"context"
	"fmt"
)

type baseStep struct {
	c *Component
}

func (s baseStep) Validate(ctx context.Context) error {
	if _, err := s.c.GetAllHosts(ctx); err != nil {
		return fmt.Errorf("MCR Step validation fail: Host validation: %s", err.Error())
	}

	return nil
}
