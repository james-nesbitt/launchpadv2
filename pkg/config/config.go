package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/cluster"
)

// ConfigBase interpretation for what type of config is being handled
type ConfigBase struct {
	APIVersion string     `yaml:"apiVersion" validate:"eq=launchpad.mirantis.com/v2.0"`
	Kind       string     `yaml:"kind" validate:"eq=cluster"`
	Metadata   ConfigMeta `yaml:"metadata"`
	Spec       yaml.Node  `yaml:"spec"`
}

// ConfigMeta defines cluster metadata.
type ConfigMeta struct {
	Name string `yaml:"name" validate:"required"`
}

// ConfigFromYamllBytes convert bytes of yaml to a cluster object
func ConfigFromYamllBytes(b []byte) (cluster.Cluster, error) {
	var cl cluster.Cluster
	var cb ConfigBase

	if err := yaml.Unmarshal(b, &cb); err != nil {
		return cl, err
	}

	if err := DecodeSpec(cb.APIVersion, &cl, cb.Spec.Decode); err != nil {
		return cl, fmt.Errorf("cluster spec decoding failed: %w", err)
	}

	return cl, nil
}
