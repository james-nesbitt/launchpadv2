package mke3

import (
	"fmt"

	dockertypesimage "github.com/docker/docker/api/types/image"

	"github.com/Mirantis/launchpad/pkg/product/common"
)

var (
	MKEProductImageRepo = "docker.io/mirantis"

	MKEBootstrapperImager = "ucp"
)

// Config has all the bits needed to configure mke during installation.
type Config struct {
	Version          string `yaml:"version" validate:"required"`
	ImageRepo        string `yaml:"imageRepo,omitempty"`
	Pull             string `yaml:"pull"`
	RegistryUserName string `yaml:"register_username"`
	RegistryPassword string `yaml:"register_password"`

	AdminUsername string `yaml:"adminUsername,omitempty"`
	AdminPassword string `yaml:"adminPassword,omitempty"`
	San           string `yaml:"san,omitempty"`

	SwarmOnly    bool         `yaml:"swarm_only"`
	InstallFlags common.Flags `yaml:"installFlags,omitempty,flow"`
	UpgradeFlags common.Flags `yaml:"upgradeFlags,omitempty,flow"`

	ConfigFile      string `yaml:"configFile,omitempty" validate:"omitempty,file"`
	ConfigData      string `yaml:"configData,omitempty"`
	LicenseFilePath string `yaml:"licenseFilePath,omitempty" validate:"omitempty,file"`

	CACertPath string `yaml:"caCertPath,omitempty" validate:"omitempty,file"`
	CertPath   string `yaml:"certPath,omitempty" validate:"omitempty,file"`
	KeyPath    string `yaml:"keyPath,omitempty" validate:"omitempty,file"`
	CACertData string `yaml:"caCertData,omitempty"`
	CertData   string `yaml:"certData,omitempty"`
	KeyData    string `yaml:"keyData,omitempty"`

	NodesHealthRetry uint `yaml:"nodesHealthRetry,omitempty" default:"0"`
}

// Prepulling images is only required if not from the main repo
func (c Config) imagePrepullRequired() bool {
	return c.ImageRepo != MKEProductImageRepo
}

// SwarmOnly return if this instance will avoid installing any kubernetes components
func (c Config) isSwarmOnly() bool {
	return c.SwarmOnly
}

// ImagePullOptions create docker image pull options for any image pulls that are needed
func (c Config) imagePullOptions() dockertypesimage.PullOptions {
	return dockertypesimage.PullOptions{}
}

func (c Config) bootstrapperImage() string {
	return fmt.Sprintf("%s/%s:%s", c.ImageRepo, MKEBootstrapperImager, c.Version)
}

func (c Config) bootStrapperInstallArgs() []string {
	ss := []string{
		"install",
		fmt.Sprintf("--admin-username='%s'", c.AdminUsername),
		fmt.Sprintf("--admin-password='%s'", c.AdminPassword),
	}

	if c.Pull != "" {
		ss = append(ss, fmt.Sprintf("--pull='%s'", c.Pull))
	}
	if c.San != "" {
		ss = append(ss, fmt.Sprintf("--san='%s'", c.San))
	}
	if c.SwarmOnly {
		ss = append(ss, "--swarm-only")
	}
	if c.RegistryUserName != "" {
		ss = append(ss, fmt.Sprintf("--registry-username='%s'", c.RegistryUserName))
		ss = append(ss, fmt.Sprintf("--registry-password='%s'", c.RegistryPassword))
	}

	return ss
}

func (c Config) bootStrapperUninstallArgs() []string {
	ss := []string{}

	return ss
}

func (c Config) bootStrapperUpgradeArgs() []string {
	ss := []string{
		"upgrade",
	}

	return ss
}
