package righost

// HostConfig configuration for an individual host.
type Config struct {
	Id    string   `yaml:"id"`
	Roles []string `yaml:"roles"`
}
