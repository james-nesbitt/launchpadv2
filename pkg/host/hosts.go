package host

// Hosts collection of Host objects for a cluster
type Hosts []*Host

// RoleHosts get all hosts with a matching role
func (hs Hosts) RoleHosts(role string) (Hosts, error) {
	rhs := []*Host{}

	for _, h := range hs {
		if h.HasRole(role) {
			rhs = append(rhs, h)
		}
	}

	return rhs, nil
}
