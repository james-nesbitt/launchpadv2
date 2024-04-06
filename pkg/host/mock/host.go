package mock

import (
	"context"
	"fmt"
	"io"

	"github.com/Mirantis/launchpad/pkg/host"
)

// ---  host ---

func NewHost(id string, roles []string, decodeerror error) host.Host {
	return ho{
		id:          id,
		Roles:       roles,
		DecodeError: decodeerror,
	}
}

// Host a host type that does nothing, but can be used as a substitute when needed
type ho struct {
	id          string
	Roles       []string `yaml:"roles"`
	DecodeError error
}

func (h ho) Id() string {
	return h.id
}

func (h ho) Exec(ctx context.Context, cmd string, inr io.Reader) (string, string, error) {
	return "see err", "Mock Host was asked to exec", fmt.Errorf("Mock Host don't exec: %+v", h)
}

func (h ho) HasRole(rs string) bool {
	for _, r := range h.Roles {
		if r == rs {
			return true
		}
	}
	return false
}
