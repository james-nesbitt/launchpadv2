package k0s

import (
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/k0sproject/version"
)

// State information for an instance of the component.
type State struct {
	// Leader is the elected leader controller host, set after discover.
	Leader *host.Host

	// HostInfo holds per-host discovery data keyed by host ID.
	HostInfo map[string]HostDiscovery

	// ClusterID is the cluster identifier discovered from a running controller, if any.
	ClusterID string
}

// HostDiscovery holds information discovered about a single host during the discover step.
type HostDiscovery struct {
	// InstalledVersion is the version of the k0s binary on disk (may differ from RunningVersion).
	InstalledVersion *version.Version

	// RunningVersion is the version of the k0s process currently running.
	// Nil if k0s is not running on this host.
	RunningVersion *version.Version

	// Role is the k0s role reported by the running process ("controller" or "worker").
	// Empty if k0s is not running.
	Role string

	// Running is true if k0s is actively running on this host.
	Running bool
}
