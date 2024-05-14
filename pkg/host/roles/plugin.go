package roles

import (
	"github.com/Mirantis/launchpad/pkg/host"
)

const (
	HostRoleRole = "role"
)

func HostGetRoleHandler(h *host.Host) HostRoleHandler {
	hge := h.MatchPlugin(HostRoleRole)
	if hge == nil {
		return nil
	}

	he, ok := hge.(HostRoleHandler)
	if !ok {
		return nil
	}

	return he
}

type HostRoleHandler interface {
	HasRole(string) bool
}

func NewHostRolesPlugin(roles []string) hostRoles {
	return hostRoles{
		Roles: roles,
	}
}

type hostRoles struct {
	Roles []string `yaml:"roles"`
}

func (hr hostRoles) Id() string {
	return "roles"
}

func (hr hostRoles) Validate() error {
	return nil
}

func (mhc hostRoles) RoleMatch(role string) bool {
	switch role {
	case HostRoleRole:
		return true
	}

	return false
}

func (hr hostRoles) HasRole(r string) bool {
	for _, hr := range hr.Roles {
		if hr == r {
			return true
		}
	}
	return false
}
