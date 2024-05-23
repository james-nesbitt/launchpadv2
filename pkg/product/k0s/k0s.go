package k0s

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/product/k0s/implementation"
)

// ValidateK0sDependencyConfig validate a K0s API client request configuration.
func (p Component) ValidateK0sDependencyConfig(dc implementation.K0sDependencyConfig) error {
	return nil
}

func (p *Component) k0sImplementation(_ context.Context) (*implementation.API, error) {
	c := implementation.Config{}
	k := implementation.NewAPI(c)
	return &k, nil
}
