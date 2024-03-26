package implementation

import "context"

// ProvidesImplementation indicates that it can create a API ImplementationAPI
type ProvidesImplementation interface {
	ImplementMKE3API(context.Context) (API, error)
}

// NewAPI Implementation
func NewAPI() API {
	return API{}
}

// API host level API implementation
type API struct {
}
