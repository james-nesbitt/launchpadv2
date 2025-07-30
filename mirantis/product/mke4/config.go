package mke4

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	kubeapi_unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kubeapi_yaml "k8s.io/apimachinery/pkg/util/yaml"
)

// Config has all the bits needed to configure mke during installation.
type Config struct {
	Name      string `yaml:"name" default:"mke4"`
	Namespace string `yaml:"namespace" default:"mke4"`

	Blueprint string `yaml:"blueprint"` // @TODO this is a placeholder, do this better with imports
}

func operatorYaml() ([]byte, error) {
	// @TODO import https://github.com/MirantisContainers/blueprint-cli/blob/main/pkg/utils/uri.go#L15
	//       instead of rewriting it here

	oryr, oryrerr := http.Get(ManifestUrlLatest)
	if oryrerr != nil {
		return []byte{}, oryrerr
	}
	defer oryr.Body.Close()

	oryrb, oryrberr := io.ReadAll(oryr.Body)
	return oryrb, oryrberr
}

// OperatorResources build the set of kubernetes resources needed to install the operator from the config.
func (c Config) OperatorResources() ([]kubeapi_unstructured.Unstructured, error) {
	var res []kubeapi_unstructured.Unstructured

	oryrb, oryrberr := operatorYaml()
	if oryrberr != nil {
		return res, oryrberr
	}

	slog.Debug("blueprint operator resource", slog.String("yaml", string(oryrb)))

	decoder := kubeapi_yaml.NewYAMLToJSONDecoder(bytes.NewReader(oryrb))

	var u kubeapi_unstructured.Unstructured
	for {
		if err := decoder.Decode(&u); err != nil {
			if err != io.EOF {
				return res, fmt.Errorf("error decoding yaml manifest file: %s", err)
			}
			break
		}
		res = append(res, u)

	}

	return res, nil
}

// BlueprintResources provide an unstructured object describing the blueprint CR from config
// @TODO this thing sucks, and should be done using structs from the Blueprint Operator.
func (c Config) BlueprintResource() (kubeapi_unstructured.Unstructured, error) {
	var bp kubeapi_unstructured.Unstructured

	bwk := fmt.Sprintf("kind: Blueprint\n%s", c.Blueprint) // the decode fails if `kind` is missing
	slog.Debug("blueprint definition", slog.String("blueprint", bwk))

	decoder := kubeapi_yaml.NewYAMLToJSONDecoder(bytes.NewReader([]byte(bwk)))

	if err := decoder.Decode(&bp); err != nil {
		if err != io.EOF {
			return bp, fmt.Errorf("error decoding yaml manifest file: %s", err)
		}
	}
	bp.SetAPIVersion("blueprint.mirantis.com/v1alpha1")
	bp.SetKind("Blueprint")
	bp.SetName(c.Name)
	bp.SetNamespace(c.Namespace)

	slog.Debug("blueprint", slog.Any("resource", bp))

	return bp, nil
}
