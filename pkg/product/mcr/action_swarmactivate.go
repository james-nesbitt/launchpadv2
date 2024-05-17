package mcr

import (
	"context"
	"fmt"
	"log/slog"

	dockertypes "github.com/docker/docker/api/types"
	dockertypesswarm "github.com/docker/docker/api/types/swarm"
	dockertypessystem "github.com/docker/docker/api/types/system"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/host/network"
	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
	"github.com/Mirantis/launchpad/pkg/implementation/docker/swarm"
)

type swarmActivateStep struct {
	baseStep

	id string
}

func (s swarmActivateStep) Id() string {
	return fmt.Sprintf("%s:mcr-swarm-activate", s.id)
}

func (s *swarmActivateStep) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "running MCR swarm-activate step", slog.String("ID", s.Id()))

	mhs, mhsgerr := s.c.GetManagerHosts(ctx)
	if mhsgerr != nil {
		return fmt.Errorf("could not retrieve managers to join the swarm: %s", mhsgerr.Error())
	}

	// 1. Find a swarm project OR Create a new one

	var l *host.Host

	if dl, err := discoverLeader(ctx, mhs); err == nil {
		slog.InfoContext(ctx, fmt.Sprintf("%s: discovered as state leader", dl.Id()), slog.Any("leader", dl))
		l = dl
	} else if il, err := initSwarm(ctx, mhs); err == nil {
		slog.InfoContext(ctx, fmt.Sprintf("%s: became new state leader from swarm init", il.Id()), slog.Any("leader", il))
		l = il
	} else {
		return fmt.Errorf("could not initialize swarm")
	}

	ld := dockerhost.HostGetDockerExec(l)

	li, lierr := ld.Info(ctx)
	if lierr != nil {
		return fmt.Errorf("%s: swarm join failed because leader docker info error: %s", l.Id(), lierr.Error())
	}
	si, sierr := ld.SwarmInspect(ctx)
	if sierr != nil {
		return fmt.Errorf("%s: swarm join failed because leader docker swarm inspect error: %s", l.Id(), sierr.Error())
	}
	ni, nierr := ld.NodeList(ctx, dockertypes.NodeListOptions{})
	if nierr != nil {
		return fmt.Errorf("%s: swarm join failed because leader docker swarm node list error: %s", l.Id(), nierr.Error())
	}

	// at this point the state should have swarm info populated
	// so all we need to do is to join the rest of the hosts
	if err := mhs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: swarm manager join", h.Id()), slog.Any("host", h))
		return joinSwarm(ctx, h, li, si, ni, "manager")
	}); err != nil {
		return fmt.Errorf("error joining managers to the swarm: %s", err.Error())
	}

	whs, whsgerr := s.c.GetWorkerHosts(ctx)
	if whsgerr != nil {
		return fmt.Errorf("could not retrieve workers to join the swarm: %s", whsgerr.Error())
	}
	if err := whs.Each(ctx, func(ctx context.Context, h *host.Host) error {
		slog.InfoContext(ctx, fmt.Sprintf("%s: swarm worker join", h.Id()), slog.Any("host", h))
		return joinSwarm(ctx, h, li, si, ni, "worker")
	}); err != nil {
		return fmt.Errorf("error joining workers to the swarm: %s", err.Error())
	}

	return nil
}

func discoverLeader(ctx context.Context, mhs host.Hosts) (*host.Host, error) {
	for _, h := range mhs {
		i, ierr := dockerhost.HostGetDockerExec(h).Info(ctx)
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

func initSwarm(ctx context.Context, mhs host.Hosts) (*host.Host, error) {
	for _, h := range mhs {
		n, nerr := network.HostGetNetwork(h).Network(ctx)
		if nerr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: could not retrieve network for swarm initialization: %s", h.Id(), nerr.Error()))
			continue
		}

		ir := dockertypesswarm.InitRequest{
			AdvertiseAddr: n.PrivateAddress,
		}

		if err := dockerhost.HostGetDockerExec(h).SwarmInit(ctx, ir); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("%s: swarm initialize fail: %s", h.Id(), err.Error()))
			continue
		}

		// swarm init suceeded on at least one manager, so we have a leader
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
func joinSwarm(ctx context.Context, h *host.Host, li dockertypessystem.Info, si dockertypesswarm.Swarm, ni []dockertypesswarm.Node, role string) error {
	// @NOTE The following three docker commands are repeated on the leader on every swarm join (probably could be cached)

	hd := dockerhost.HostGetDockerExec(h) // we already tested for nil when we build the host lists

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
