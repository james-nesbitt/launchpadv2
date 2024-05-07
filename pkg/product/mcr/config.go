package mcr

// Config definition for the MCR product.
type Config struct {
	Version           string `yaml:"version" json:"version"`
	RepoURL           string `yaml:"repoURL,omitempty" json:"repoURL,omitempty"`
	InstallURLLinux   string `yaml:"installURLLinux,omitempty" json:"installURLLinux,omitempty"`
	InstallURLWindows string `yaml:"installURLWindows,omitempty" json:"installURLWindows,omitempty"`
	Channel           string `yaml:"channel,omitempty" json:"channel,omitempty"`

	SwarmInstallFlags []string `yaml:"swarm_install_flags,omitempty" json:"swarm_init_flags,omitempty"`
}
