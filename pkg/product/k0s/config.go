package k0s

// Config definition for the MCR product.
type Config struct {
	Instance string `yaml:"instance" json:"instance"`
	Version  string `yaml:"version" json:"version"`
}
