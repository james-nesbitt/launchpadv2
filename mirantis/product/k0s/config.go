package k0s

import (
	"github.com/k0sproject/version"
)

var (
	K0sReleaseLinkBase = "https://github.com/k0sproject/k0s/releases/download"
)

// Config definition for the MCR product.
type Config struct {
	Version        version.Version `yaml:"version" json:"version"`
	VersionChannel string          `yaml:"versionChannel,omitempty"`
	DynamicConfig  bool            `yaml:"dynamicConfig,omitempty" default:"false"`
	K0sConfig      K0sConfig       `yaml:"config,omitempty"`
	Metadata       K0sMetadata     `yaml:"-"`
}

// K0sMetadata contains gathered information about k0s cluster.
type K0sMetadata struct {
	ClusterID        string
	VersionDefaulted bool
}

// ShouldDownload should we download binaries directly to the host, os should we upload from the running client machine.
func (c Config) ShouldDownload() bool {
	return true
}
