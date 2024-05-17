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
