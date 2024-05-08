package mcr

import (
	"sync"

	dockertypes "github.com/docker/docker/api/types"
	dockertypessystem "github.com/docker/docker/api/types/system"

	dockerhost "github.com/Mirantis/launchpad/pkg/implementation/docker/host"
)

// state state object for MCR, to be used to compute delta from config.
type state struct {
	leader *dockerhost.Host

	hosts hostStates
}

type hostStates map[string]*hostState

var (
	hostsLock = sync.Mutex{}
)

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

func (hs *hostState) setVersion(v map[string]dockertypes.Version) {
	hs.mu.Lock()
	hs.version = v
	hs.mu.Unlock()
}
