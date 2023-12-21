package mke

// NewMKE MKE constructor.
func NewMKE(c Config) MKE {
	return MKE{config: c}
}

// MKE product implementation.
type MKE struct {
	config Config
}

// Debug product debug.
func (m MKE) Debug() interface{} {
	return nil
}
