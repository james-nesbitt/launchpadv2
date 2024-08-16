package msr4

import (
	"context"
	"fmt"

	"helm.sh/helm/v3/pkg/release"
)

func (c *Component) getRelease(ctx context.Context) (*release.Release, error) {
	hc, hcerr := c.helmClient(ctx)
	if hcerr != nil {
		return nil, hcerr
	}

	r, rerr := hc.GetRelease(c.Name())
	if rerr != nil {
		return nil, rerr
	}

	return r, nil
}

func (c *Component) addMsr4HelmRepo(ctx context.Context) error {
	hc, hcerr := c.helmClient(ctx)
	if hcerr != nil {
		return hcerr
	}

	hr := c.helmRepo()

	if err := hc.AddOrUpdateChartRepo(hr); err != nil {
		return fmt.Errorf("Error adding repo: %s", err.Error())
	}

	return nil
}

func (c *Component) installMsr4Release(ctx context.Context) (*release.Release, error) {
	hc, hcerr := c.helmClient(ctx)
	if hcerr != nil {
		return nil, hcerr
	}

	hcs := c.helmChartSpec()
	hcopts := c.helmChartOpts()

	r, rerr := hc.InstallOrUpgradeChart(ctx, &hcs, &hcopts)
	if rerr != nil {
		return nil, rerr
	}

	return r, nil
}

func (c *Component) uninstallMsr4Release(ctx context.Context) error {
	hc, hcerr := c.helmClient(ctx)
	if hcerr != nil {
		return hcerr
	}

	if err := hc.UninstallReleaseByName(c.Name()); err != nil {
		return err
	}

	return nil
}
