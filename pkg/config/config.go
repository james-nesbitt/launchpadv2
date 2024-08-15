package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/Mirantis/launchpad/pkg/project"
	"github.com/creasty/defaults"
)

// ConfigBase interpretation for what type of config is being handled.
type ConfigBase struct {
	APIVersion string     `yaml:"apiVersion" validate:"eq=launchpad.mirantis.com/v2.0"`
	Kind       string     `yaml:"kind" validate:"eq=project"`
	Metadata   ConfigMeta `yaml:"metadata"`
	Spec       yaml.Node  `yaml:"spec" validate:"required"`
}

// ConfigMeta defines project metadata.
type ConfigMeta struct {
	Name string `yaml:"name" validate:"required"`
}

// ConfigFromYamllBytes convert bytes of yaml to a project object.
func ConfigFromYamllBytes(b []byte) (project.Project, error) {
	var cl project.Project = project.New()
	var cb ConfigBase

	defaults.Set(&cb)

	if err := yaml.Unmarshal(b, &cb); err != nil {
		return cl, err
	}

	if err := DecodeSpec(cb.APIVersion, &cl, cb.Spec.Decode); err != nil {
		return cl, fmt.Errorf("project spec decoding failed: %w", err)
	}

	return cl, nil
}
