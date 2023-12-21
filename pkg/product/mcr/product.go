package mcr

// NewMCR constructor for MCR from config.
func NewMCR(c Config) MCR {
	return MCR{config: c}
}

// MCR product implementation.
type MCR struct {
	config Config
}

// Debug product debug.
func (m MCR) Debug() interface{} {
	return nil
}

// State discover the state of the cluster installaton
