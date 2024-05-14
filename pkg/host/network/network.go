package network

import (
	"context"

	"github.com/Mirantis/launchpad/pkg/host"
)

const (
	HostRoleNetwork = "network"
)

func HostGetExecutor(h *host.Host) HostNetwork {
	hgn := h.MatchPlugin(HostRoleNetwork)
	if hgn == nil {
		return nil
	}

	hn, ok := hgn.(HostNetwork)
	if !ok {
		return nil
	}

	return hn
}

type HostNetwork interface {
	Network(ctx context.Context) Network
}

// Network networking information about the Host
type Network struct {
	PublicAddress    string `yaml:"public_address"`
	PrivateInterface string `yaml:"private_interface"`
	PrivateAddress   string `yaml:"private_address"`
}
