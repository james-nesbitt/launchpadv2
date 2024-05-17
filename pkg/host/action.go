package host

import (
	"context"
	"errors"
)

type baseStep struct {
	c *HostsComponent
}

func (s baseStep) Validate(ctx context.Context) error {
	if err := s.c.Validate(ctx); err != nil {
		return err
	}

	errs := []error{}
	for _, d := range s.c.deps {
		if err := d.Validate(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
