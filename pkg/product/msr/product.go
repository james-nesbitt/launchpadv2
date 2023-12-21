package msr

// NewMKE MKE constructor.
func NewMSR(c Config) MSR {
	return MSR{config: c}
}

// MSR product implementation.
type MSR struct {
	config Config
}

// Debug product debug.
func (m MSR) Debug() interface{} {
	return nil
}
