package msr2

import (
	"context"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/dependency"
	"github.com/Mirantis/launchpad/pkg/product/mke3"
)

func (p *MSR2) getMK3(ctx context.Context) (*mke3.MKE3, error) {
	if p.mke3r == nil {
		return nil, fmt.Errorf("%s: %s dependency missing", ComponentType, mke3.ComponentType)
	}

	d := p.mke3r.Matched(ctx)
	if d == nil {
		return nil, fmt.Errorf("%s: %s dependency not Matched", ComponentType, mke3.ComponentType)
	}

	if err := d.Met(ctx); err != nil {
		return nil, fmt.Errorf("%s: %s dependency not Met: %w", ComponentType, mke3.ComponentType, err)
	}

	mke3d, ok := d.(mke3.MKE3Dependency)
	if !ok {
		return nil, fmt.Errorf("%s: %s dependency wrong type: %+v", ComponentType, mke3.ComponentType, d)
	}

	mke3c := mke3d.ProvidesMKE3(ctx)
	if mke3c == nil {
		return nil, fmt.Errorf("%s: %s dependency gave empty client: %+v", ComponentType, mke3.ComponentType, mke3d)
	}

	return mke3c, nil

}

// Requires declare that we need a HostsRoles dependency
func (p *MSR2) Requires(_ context.Context) dependency.Requirements {
	p.mke3r = mke3.NewMKE3Requirement(
		p.Name(),
		fmt.Sprintf("%s: relies on MKE for enzi", ComponentType),
		mke3.MKE3DependencyConfig{},
	)

	return dependency.Requirements{p.mke3r}
}
