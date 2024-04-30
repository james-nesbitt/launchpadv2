package kubernetes

const (
	ImplementationType = "kubernetes"
)

type Version struct {
	Version string
	APIs    []string
}

// NewKubernetes constructor
func NewKubernetes(config Config) *Kubernetes {
	return &Kubernetes{
		config: config,
	}
}

// Kubernetes implementation
type Kubernetes struct {
	config Config
}

// IsValidaKubernetesVersion is the k8s version an acceptable value for this code base
func IsValidKubernetesVersion(_ Version) error {
	return nil
}
