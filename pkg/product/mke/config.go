package mke

import (
	"github.com/Mirantis/launchpad/pkg/product/common"
)

// Config has all the bits needed to configure mke during installation.
type Config struct {
	Version             string       `yaml:"version" validate:"required"`
	ImageRepo           string       `yaml:"imageRepo,omitempty"`
	AdminUsername       string       `yaml:"adminUsername,omitempty"`
	AdminPassword       string       `yaml:"adminPassword,omitempty"`
	InstallFlags        common.Flags `yaml:"installFlags,omitempty,flow"`
	SwarmInstallFlags   common.Flags `yaml:"swarmInstallFlags,omitempty,flow"`
	SwarmUpdateCommands []string     `yaml:"swarmUpdateCommands,omitempty,flow"`
	UpgradeFlags        common.Flags `yaml:"upgradeFlags,omitempty,flow"`
	ConfigFile          string       `yaml:"configFile,omitempty" validate:"omitempty,file"`
	ConfigData          string       `yaml:"configData,omitempty"`
	LicenseFilePath     string       `yaml:"licenseFilePath,omitempty" validate:"omitempty,file"`
	CACertPath          string       `yaml:"caCertPath,omitempty" validate:"omitempty,file"`
	CertPath            string       `yaml:"certPath,omitempty" validate:"omitempty,file"`
	KeyPath             string       `yaml:"keyPath,omitempty" validate:"omitempty,file"`
	CACertData          string       `yaml:"caCertData,omitempty"`
	CertData            string       `yaml:"certData,omitempty"`
	KeyData             string       `yaml:"keyData,omitempty"`
	NodesHealthRetry    uint         `yaml:"nodesHealthRetry,omitempty" default:"0"`
}
