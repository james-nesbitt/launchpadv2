package msr4

import (
	"context"
)

type baseStep struct {
	c *Component
}

func (s baseStep) Validate(ctx context.Context) error {
	return nil
}
