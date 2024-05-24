package host

import (
	"context"
	"errors"
)

var (
	ErrNoHostPluginHandlerRegistered = errors.New("No host plugin handlers have been registered")
	ErrUnknownHostPluginType         = errors.New("Unknown host plugin handler")

	// hostPluginFactories the global map of PluginFactories, managed by
	// the register function below
	hostPluginFactories = map[string]HostPluginFactory{}
)

// RegisterHostPlugin register a host plugin factory
func RegisterHostPluginFactory(key string, hp HostPluginFactory) {
	hostPluginFactories[key] = hp
}

// HostPluginFactory a factory for host plugins injected into the host package
//
// The factories must be able to create new plugins (usually from decoding)
// but also server additional roles as mutators, for things like injecting
// cli commands etc. To perform these additional roles, the PluginFactory
// should also implement other interfaces.
type HostPluginFactory interface {
	// HostPlugin build a new host plugin
	HostPlugin(context.Context, *Host) HostPlugin
	// HostPluginDecode provide a Host Plugin decoder function
	//
	// The decoder function is ugly, but it is meant to to take a
	// yaml/json .HostPluginDecode() function, and turn it into a plugin
	HostPluginDecode(context.Context, *Host, func(interface{}) error) (HostPlugin, error)
}

// HostPlugin the base interface for host plugins
//
// Note that plugins are mutable, and the other interfaces
// are more important than this interface
type HostPlugin interface {
	// Id return the host unique identifier
	Id() string
	// RoleMatch find a host plugin that matches a role
	RoleMatch(role string) bool
	// Validate is the host plugin valid after configuraiton
	Validate() error
}

// MatchPlugin finc the host plugin that fullfills a role
func (h *Host) MatchPlugin(role string) interface{} {
	for _, hc := range h.plugins {
		if hc.RoleMatch(role) {
			return hc
		}
	}
	return nil
}

// PluginIds retrieve the IDs of all of the plugin for the host
func (h *Host) PluginIDs() []string {
	pids := []string{}
	for _, p := range h.plugins {
		pids = append(pids, p.Id())
	}
	return pids
}

// HasPlugin does the host have a plugin of the passed id
func (h *Host) HasPlugin(pid string) bool {
	for _, p := range h.plugins {
		if pid == p.Id() {
			return true
		}
	}
	return false
}

// AddPlugin add a new pluging to the host
func (h *Host) AddPlugin(p HostPlugin) {
	h.plugins = append(h.plugins, p)
}
