package implementation

type Config struct {
}

// NewAPI Implementation.
func NewAPI(c Config) API {
	return API{
		config: c,
	}
}

// API host level API implementation.
type API struct {
	config Config
}
