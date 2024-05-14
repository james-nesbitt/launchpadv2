package host

func NewHost(id string) *Host {
	return &Host{
		id: id,
	}
}

type Host struct {
	id      string
	plugins []HostPlugin
}

func (h Host) Id() string {
	return h.id
}

func (h *Host) MatchPlugin(role string) interface{} {
	for _, hc := range h.plugins {
		if hc.RoleMatch(role) {
			return hc
		}
	}
	return nil
}

func (h *Host) PluginIDs() []string {
	pids := []string{}
	for _, p := range h.plugins {
		pids = append(pids, p.Id())
	}
	return pids
}

func (h *Host) HasPlugin(pid string) bool {
	for _, p := range h.plugins {
		if pid == p.Id() {
			return true
		}
	}
	return false
}

func (h *Host) AddPlugin(p HostPlugin) {
	h.plugins = append(h.plugins, p)
}
