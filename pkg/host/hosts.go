package host

import (
	"errors"
	"fmt"
)

var (
	ErrHostsFilterInsufficient = errors.New("not enough matching hosts to filter")
)

// Hosts collection of Host objects for a cluster.
type Hosts map[string]Host

// HostsRolesFilter filter a Hosts set based on filter arguments.
type HostsRolesFilter struct {
	Roles []string
	Min   int
	Max   int
}

// FilterRoles filter hosts down to matching roles, enforcing a minimum and accepting a maximum count.
func (hs Hosts) FilterRoles(filter HostsRolesFilter) (Hosts, error) {
	fhs := Hosts{}

	hasRoles := func(h Host, rs []string) bool {
		if len(rs) == 0 {
			return true
		}

		for _, r := range rs {
			if h.HasRole(r) {
				return true
			}
		}
		return false
	}

	for id, h := range hs {
		if hasRoles(h, filter.Roles) {
			fhs[id] = h
		}
		if filter.Max > 0 && len(fhs) >= filter.Max {
			break
		}
	}

	if filter.Min < 0 && len(fhs) < filter.Min {
		return fhs, fmt.Errorf("%w: %+v", ErrHostsFilterInsufficient, filter)
	}

	return fhs, nil
}
