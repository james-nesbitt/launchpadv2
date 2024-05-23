package project

const (
	// Apply command means ensure that resources are in place.
	CommandKeyApply = "apply"
	// Reset command means ensure that resources are removed.
	CommandKeyReset = "reset"
	// Discover command means discover what resources are in place (state).
	CommandKeyDiscover = "discover"
)
