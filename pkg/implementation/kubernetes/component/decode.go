package kubernetes

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/creasty/defaults"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/implementation/kubernetes"
)

func init() {
	component.RegisterDecoder(ComponentType, DecodeComponent)
}

// DecodeComponent decode a new component from an unmarshall decoder.
func DecodeComponent(id string, d func(interface{}) error) (component.Component, error) {
	var mc mockConfig
	defaults.Set(&mc)

	if err := d(&mc); err != nil {
		return nil, fmt.Errorf("Failure to decode product '%s' : %w", ComponentType, err)
	}

	var c kubernetes.Config

	if mc.KubeConfigFile != "" {
		fr, frerr := os.Open(mc.KubeConfigFile)
		if frerr != nil {
			return nil, fmt.Errorf("%s couldn't open kube config file %s : %s", id, mc.KubeConfigFile, frerr.Error())
		}
		defer fr.Close()

		frb, frberr := io.ReadAll(fr)
		if frberr != nil {
			return nil, fmt.Errorf("%s couldn't open kube config file %s : %s", id, mc.KubeConfigFile, frberr.Error())
		}

		frbc, frbcerr := kubernetes.ConfigFromKubeConf(frb)
		if frbcerr != nil {
			return nil, fmt.Errorf("%s couldn't convert kube config file %s to kubernetes config object : %s", id, mc.KubeConfigFile, frbcerr.Error())
		}

		c = frbc
	} else if mc.KubeConfigYaml != "" {
		cybc, cybccerr := kubernetes.ConfigFromKubeConf([]byte(mc.KubeConfigYaml))
		if cybccerr != nil {
			return nil, fmt.Errorf("%s couldn't convert kube config yaml to kubernetes config object : %s", id, cybccerr.Error())
		}

		c = cybc
	} else {
		slog.Warn("no explicit kubecontext source given (currently unsupported)")
		dc, dcerr := kubernetes.ConfigFromEnv(nil)
		if dcerr != nil {
			return nil, fmt.Errorf("%s default config load failed: %s", id, dcerr.Error())
		}

		c = dc
	}

	// @TODO find a way to allow overrides

	// @TODO Validate that config was populated

	return NewComponent(id, c), nil
}

// decoding target to determine where we should get kube config from
// @TODO figure out how to add overrides (especially for auth)
type mockConfig struct {
	KubeConfigYaml string `yaml:"kubeYaml" json:"kubeYaml"`
	KubeConfigFile string `yaml:"kubeFile" json:"kubeFile"`
}
