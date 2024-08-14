package msr4

import helmclient "github.com/mittwald/go-helm-client"

// Config for MSR.
type Config struct {
	Version   string `yaml:"version"`
	Namespace string `yaml:"namespace" default:"msr4"`

	RepositoryCache  string `yaml:"repocache" default:"/tmp/.helmcache"`
	RepositoryConfig string `yaml:"repoconfig" default:"/tmp/.helmrepo"`
	Debug            bool   `yaml:"debug"`
	Linting          bool   `yaml:"lint"`
}

func (c Config) HelmOptions() helmclient.Options {
	return helmclient.Options{
		Namespace:        c.Namespace,
		RepositoryCache:  c.RepositoryCache,
		RepositoryConfig: c.RepositoryConfig,
		Debug:            c.Debug,
		Linting:          c.Linting,
	}
}
