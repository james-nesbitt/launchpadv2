package host

// NewHost Host constructor from HostConfig
func NewHost(hc HostConfig) Host {
	return Host{config: hc}
}

// Host definition for one cluster target
type Host struct {
	config       HostConfig
	anotherthing string
}

// HasRole answer boolean if a host has a role
func (h Host) HasRole(role string) bool {
	return h.config.Role == role
}
