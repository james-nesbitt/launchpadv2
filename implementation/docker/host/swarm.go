package dockerhost

import (
	"context"
	"fmt"
	"log/slog"

	dockertypesswarm "github.com/docker/docker/api/types/swarm"
	dockertypessystem "github.com/docker/docker/api/types/system"

	"github.com/Mirantis/launchpad/implementation/docker/swarm"
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

func DiscoverLeader(ctx context.Context, mhs host.Hosts) (*host.Host, error) {
	for _, h := range mhs {
		i, ierr := HostGetDockerExec(h).Info(ctx)
		if ierr != nil {
			continue
		}

		if i.Swarm.ControlAvailable {
			slog.DebugContext(ctx, fmt.Sprintf("%s: swarm already active, this host can act as a leader.", h.Id()), slog.Any("host", h))
			return h, nil
		}
	}
	return nil, fmt.Errorf("No swarm leader discovered")
}

func InitSwarm(ctx context.Context, mhs host.Hosts) (*host.Host, error) {
	for _, h := range mhs {
		n, nerr := network.HostGetNetwork(h).Network(ctx)
		if nerr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: could not retrieve network for swarm initialization: %s", h.Id(), nerr.Error()))
			continue
		}

		ir := dockertypesswarm.InitRequest{
			AdvertiseAddr: n.PrivateAddress,
		}

		if err := HostGetDockerExec(h).SwarmInit(ctx, ir); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("%s: swarm initialize fail: %s", h.Id(), err.Error()))
			continue
		}

		// swarm init succeeded on at least one manager, so we have a leader
		slog.DebugContext(ctx, fmt.Sprintf("%s: swarm initialized. This is the new swarm leader", h.Id()), slog.Any("host", h))

		return h, nil
	}

	return nil, fmt.Errorf("could not initialize swarm on any manager host")
}

// join a host to the swarm
//
// @NOTE We have to check if the host is already in the swarm
//
//	and we may need to leave a/the swarm first
func JoinSwarm(ctx context.Context, h *host.Host, li dockertypessystem.Info, si dockertypesswarm.Swarm, ni []dockertypesswarm.Node, role string) error {
	// @NOTE The following three docker commands are repeated on the leader on every swarm join (probably could be cached)

	hd := HostGetDockerExec(h) // we already tested for nil when we build the host lists

	if hi, err := hd.Info(ctx); err == nil && hi.Swarm.LocalNodeState != "inactive" {
		if hi.Swarm.NodeID == li.Swarm.NodeID {
			slog.InfoContext(ctx, fmt.Sprintf("%s: host already in swarm as the leader", h.Id()))
			return nil
		}

		for _, sn := range ni {
			if sn.ID == hi.Swarm.NodeID {
				// Node is already in swarm

				if hi.Swarm.ControlAvailable && role == "worker" {
					slog.WarnContext(ctx, fmt.Sprintf("%s: host is already in the swarm, but is a worker when it is supposed to be a manager", h.Id()))
				} else if !hi.Swarm.ControlAvailable && role == "manager" {
					slog.WarnContext(ctx, fmt.Sprintf("%s: host is already in the swarm, but is a manager when it is supposed to be a worker", h.Id()))
				}

				// We should handle the above issues by promoting/demoting the hosts, but we are stuck in a place where
				// the docker cli doesn't give that option.  We could make the host leave the swarm, but then it could
				// disrupt the workload.

				slog.InfoContext(ctx, fmt.Sprintf("%s: host already in swarm", h.Id()), slog.Any("host-swarm", hi.Swarm))
				return nil

			}
		}

		// Node is in a swarm, but it is the wrong swarm
		if lerr := hd.SwarmLeave(ctx, true); lerr != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("%s: host already in swarm", h.Id()), slog.Any("host-swarm", hi.Swarm), slog.Any("leader-swarm", li))
			return fmt.Errorf("%s: was in another swarm, and failed to leave it", h.Id())
		}
	}

	r := dockertypesswarm.JoinRequest{
		RemoteAddrs: []string{swarm.SwarmAddress(li.Swarm.NodeAddr)},
	}

	switch role {
	case "manager":
		r.JoinToken = si.JoinTokens.Manager
	case "worker":
		r.JoinToken = si.JoinTokens.Worker
	}

	if err := hd.SwarmJoin(ctx, r); err != nil {
		return fmt.Errorf("%s: failed to join manager to swarm: %s", h.Id(), err.Error())
	}

	return nil
}
