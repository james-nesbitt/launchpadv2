package roles

import (
	"errors"
	"fmt"

	"github.com/Mirantis/launchpad/pkg/host"
)

var (
	ErrHostsFilterInsufficient = errors.New("not enough matching hosts to filter")
)

// HostsRolesFilter filter a Hosts set based on filter arguments.
type HostsRolesFilter struct {
	Roles []string
	Min   int
	Max   int
}

// FilterRoles filter hosts down to matching roles, enforcing a minimum and accepting a maximum count.
func FilterRoles(hs host.Hosts, filter HostsRolesFilter) (host.Hosts, error) {
	fhs := host.Hosts{}

	for _, h := range hs {
		if len(filter.Roles) == 0 {
			fhs.Add(h)
		} else if hgr := HostGetRoleHandler(h); hgr == nil {
			continue
		} else {
			for _, r := range filter.Roles {
				if hgr.HasRole(r) {
					fhs.Add(h)
					break
				}
			}
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
