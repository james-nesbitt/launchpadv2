package mcr

// State state object for MCR, to be used to compute delta from config.
type State struct {
	Version           string `yaml:"version" json:"version"`
	RepoURL           string `yaml:"repoURL,omitempty" json:"repoURL,omitempty"`
	InstallURLLinux   string `yaml:"installURLLinux,omitempty" json:"installURLLinux,omitempty"`
	InstallURLWindows string `yaml:"installURLWindows,omitempty" json:"installURLWindows,omitempty"`
	Channel           string `yaml:"channel,omitempty" json:"channel,omitempty"`
}
