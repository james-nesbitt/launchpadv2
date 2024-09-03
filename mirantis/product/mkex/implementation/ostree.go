package implementation

type RpmOSTreeStatusDeployment struct {
	RequestedLocalPackages            []string                              `json:"requested-local-packages"`
	BaseCommitMeta                    []RpmOSTreeStatusDeploymentCommitMeta `json:"base-commit-meta"`
	BaseRemovals                      []string                              `json:"base-removals"`
	Unlocked                          string                                `json:"unlocked"`
	Booted                            bool                                  `json:"booted"`
	RequestedLocalFileoveridePackages []string                              `json:"requested-local-fileoverride-packages"`
	RequestedModulesEnabled           []string                              `json:"requested-modules-enabled"`
	ID                                string                                `json:"id"`
	BaseRemoteReplacements            map[string]string                     `json:"base-remote-replacements"`
	OSName                            string                                `json:"osname"`
	Pinned                            bool                                  `json:"pinned"`
	BaseLocalReplacements             []string                              `json:"base-local-replacements"`
	RegenarateInitramfs               bool                                  `json:"regenerate-initramfs"`
	Origin                            string                                `json:"origin"`
	Packages                          []string                              `json:"packages"`
	Version                           string                                `json:"version"`
	Checksum                          string                                `json:"checksum"`
	RequestedBaseLocalReplacements    []string                              `json:"requested-base-local-replacements"`
	RequestedModules                  []string                              `json:"requested-modules"`
	RequestedPackages                 []string                              `json:"requested-packages"`
	Serial                            int                                   `json:"serial"`
	Timestamp                         string                                `json:"timestamp"`
	Staged                            bool                                  `json:"staged"`
	RequestedBaseRemovals             []string                              `json:"requested-base-removals"`
	GbgEnabled                        bool                                  `json:"gpg-enabled"`
	Modules                           []string                              `json:"modules"`
}

type RpmOSTreeStatusDeploymentCommitMeta struct {
	OSTreeBootable          bool                            `json:"ostree.bootable"`
	RPMOStreeInputHash      string                          `json:"rpmostree.inputhash"`
	Version                 string                          `json:"version"`
	OSTreeComposeFSDigestV0 []string                        `json:"ostree.composefs.digest.v0"`
	OSTreeLinux             string                          `json:"ostree.linux"`
	RPMOSTreeRpmdRepos      []RpmOSTreeStatusDeploymentRepo `json:"rpmostree.rpmmd-repos"`
}

type RpmOSTreeStatusDeploymentRepo struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

func DecodeStatus(j string) []RpmOSTreeStatusDeployment {
	return []RpmOSTreeStatusDeployment{}
}
