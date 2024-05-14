package host

type HostPlugin interface {
	// Id return the host unique identifier
	Id() string
	// RoleMatch find a host plugin that matches a role
	RoleMatch(role string) bool
	// Validate is the host plugin valid after configuraiton
	Validate() error
}
