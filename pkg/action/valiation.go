package action

import "context"

// Validator can validate itself with just a context.
type Validator interface {
	// Validate that the phase can Run
	Validate(context.Context) error
}
