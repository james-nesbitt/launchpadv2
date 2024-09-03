package swarm

type Config struct {
	Host          string
	Key           string
	Cert          string
	CACert        string
	TLSSkipVerify bool
}
