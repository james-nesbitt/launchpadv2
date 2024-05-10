package mcr

import (
	"sync"
)

var (
	stateLock = sync.Mutex{}
)

// state state object for MCR, to be used to compute delta from config.
type state struct {
}

// Debug state for output
func (s *state) Debug() interface{} {
	return nil
}
