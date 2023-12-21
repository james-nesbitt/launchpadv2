package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/cluster"
	v2_0 "github.com/Mirantis/launchpad/pkg/config/v2_0"
)

const (
	API_VERSION_2_0 = "launchpad.mirantis.com/v2.0"
)

// ConfigBase interpretation for what type of config is being handled
type ConfigBase struct {
	APIVersion string     `yaml:"apiVersion" validate:"eq=launchpad.mirantis.com/v2.0"`
	Kind       string     `yaml:"kind" validate:"eq=cluster"`
	Metadata   ConfigMeta `yaml:"metadata"`
}

// ConfigMeta defines cluster metadata.
type ConfigMeta struct {
	Name string `yaml:"name" validate:"required"`
}

// ConfigFromYamllBytes convert bytes of yaml to a cluster object
func ConfigFromYamllBytes(b []byte) (cluster.Cluster, error) {
	var cluster cluster.Cluster
	var cb ConfigBase

	if err := yaml.Unmarshal(b, &cb); err != nil {
		return cluster, err
	}

	switch cb.APIVersion {
	case API_VERSION_2_0:
		var v20conf v2_0.Config
		if err := yaml.Unmarshal(b, &v20conf); err != nil {
			return cluster, err
		}
		if c, err := v20conf.Cluster(); err != nil {
			return cluster, err
		} else {
			cluster = c
		}

	default:
		return cluster, errors.New(fmt.Sprintf("apiVersion '%s' is not compatible", cb.APIVersion))
	}

	return cluster, nil
}
