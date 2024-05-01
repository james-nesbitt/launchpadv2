package implementation

type Config struct {
	Version string
}

// NewAPI Implementation.
func NewAPI(c Config) *API {
	return &API{
		config: c,
	}
}

// API host level API implementation.
type API struct {
	config Config
}

func (a API) Version() string {
	return a.config.Version
}
