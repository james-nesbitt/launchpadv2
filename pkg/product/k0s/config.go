package k0s

import (
	"fmt"

	"github.com/k0sproject/dig"

	"github.com/k0sproject/version"
)

var (
	K0sReleaseLinkBase = "https://github.com/k0sproject/k0s/releases/download/"
)

// Config definition for the MCR product.
type Config struct {
	Version        version.Version `yaml:"version" json:"version"`
	VersionChannel string          `yaml:"versionChannel,omitempty"`
	DynamicConfig  bool            `yaml:"dynamicConfig,omitempty" default:"false"`
	K0sConfig      dig.Mapping     `yaml:"config,omitempty"`
	Metadata       K0sMetadata     `yaml:"-"`
}

// K0sMetadata contains gathered information about k0s cluster
type K0sMetadata struct {
	ClusterID        string
	VersionDefaulted bool
}

func (c Config) DownloadURL(arch string) string {
	// https://github.com/k0sproject/k0s/releases/download/v1.30.0%2Bk0s.0/k0s-v1.30.0+k0s.0-amd64
	// v1.30.0%2Bk0s.0/k0s-v1.30.0+k0s.0-amd64

	return fmt.Sprintf("https://github.com/k0sproject/k0s/releases/download/%[1]s/k0s-%[1]s-%[2]s", c.Version.String(), arch)
}
