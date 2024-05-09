package mcr

import (
	"sync"

	dockertypes "github.com/docker/docker/api/types"
	dockertypesswarm "github.com/docker/docker/api/types/swarm"
	dockertypessystem "github.com/docker/docker/api/types/system"

	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

var (
	stateLock = sync.Mutex{}
	hostsLock = sync.Mutex{}
)

// state state object for MCR, to be used to compute delta from config.
type state struct {
	leader *dockerhost.Host
	swarm  *dockertypesswarm.Swarm

	hosts hostStates
}

func (s *state) SetLeader(l *dockerhost.Host) {
	stateLock.Lock()
	s.leader = l
	stateLock.Unlock()
}

func (s state) GetLeader() *dockerhost.Host {
	stateLock.Lock()
	defer stateLock.Unlock()
	return s.leader
}

type hostStates map[string]*hostState

func (hss hostStates) GetOrCreate(id string) *hostState {
	hostsLock.Lock()
	defer hostsLock.Unlock()

	if hs, ok := hss[id]; ok {
		return hs
	}

	hs := &hostState{
		version: map[string]dockertypes.Version{},
	}

	hss[id] = hs

	return hs
}

type hostState struct {
	mu sync.Mutex

	version map[string]dockertypes.Version
	info    dockertypessystem.Info
}

func (hs *hostState) setInfo(i dockertypessystem.Info) {
	hs.mu.Lock()
	hs.info = i
	hs.mu.Unlock()
}

func (hs *hostState) getInfo() dockertypessystem.Info {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	return hs.info
}

func (hs *hostState) setVersion(v map[string]dockertypes.Version) {
	hs.mu.Lock()
	hs.version = v
	hs.mu.Unlock()
}

func (hs *hostState) getVersion() map[string]dockertypes.Version {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	return hs.version
}
