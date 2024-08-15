package msr4

import (
	"time"

	helmc "github.com/mittwald/go-helm-client"
	helmrepo "helm.sh/helm/v3/pkg/repo"
)

func (c Component) helmOptions() helmc.Options {
	return c.config.HelmOptions()
}

func (c Component) helmRepo() helmrepo.Entry {
	return helmrepo.Entry{
		Name: "msr4",
		URL:  "https://registry.mirantis.com/charts/harbor/helm",
	}
}

func (c Component) helmChartSpec() helmc.ChartSpec {
	cs := helmc.ChartSpec{
		ChartName:   "msr4/harbor",
		ReleaseName: c.Name(),
		Version:     c.config.Version,
		Namespace:   c.config.Namespace,

		CreateNamespace: true,
		UpgradeCRDs:     true,
		Wait:            true,
		WaitForJobs:     true,
		Timeout:         time.Second * 300,

		ValuesYaml: `
persistence:
  enabled: False
`,
	}

	return cs
}

func (c Component) helmChartOpts() helmc.GenericHelmOptions {
	return helmc.GenericHelmOptions{}
}
