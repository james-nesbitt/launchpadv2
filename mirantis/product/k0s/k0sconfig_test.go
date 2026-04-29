package k0s_test

import (
	"testing"

	"github.com/Mirantis/launchpad/mirantis/product/k0s"
	"gopkg.in/yaml.v3"
)

func Test_K0sConfigUnMarshall(t *testing.T) {
	k0sConfigYaml := `
apiVersion: k0s.k0sproject.io/v1beta1
kind: ClusterConfig
metadata:
  name: k0s
spec:
  api:
    externalAddress: my.lb.org
    k0sApiPort: 9443
    port: 6443
    sans:
    - my.san.org`

	var kc k0s.K0sConfig
	if err := yaml.Unmarshal([]byte(k0sConfigYaml), &kc); err != nil {
		t.Errorf("failed to unmarshal: %s", err.Error())
	}

	if kc.DigString("spec", "api", "externalAddress") != "my.lb.org" {
		t.Errorf("wrong spec.api.externalAddress: %s", kc.DigString("spec", "api", "externalAddress"))
	}

	apiMap := kc.DigMapping("spec", "api")
	sans, ok := apiMap["sans"].([]interface{})
	if !ok || len(sans) == 0 {
		t.Fatalf("expected spec.api.sans to have at least one entry")
	}
	s0, ok := sans[0].(string)
	if !ok {
		t.Fatalf("spec.api.sans[0] is not a string: %T", sans[0])
	}
	if s0 != "my.san.org" {
		t.Errorf("wrong spec.api.sans[0]: %v", sans[0])
	}
}

func Test_K0sConfigRoundTrip(t *testing.T) {
	k0sConfigYaml := `
apiVersion: k0s.k0sproject.io/v1beta1
kind: ClusterConfig
metadata:
  name: k0s
spec:
  api:
    address: 10.0.0.1
    port: 6443`

	var kc k0s.K0sConfig
	if err := yaml.Unmarshal([]byte(k0sConfigYaml), &kc); err != nil {
		t.Fatalf("failed to unmarshal: %s", err.Error())
	}

	// Mutate via dig
	kc.DigMapping("spec", "api")["address"] = "192.168.1.1"

	// Marshal back to YAML
	out, err := yaml.Marshal(kc)
	if err != nil {
		t.Fatalf("failed to marshal: %s", err.Error())
	}

	// Unmarshal and verify
	var kc2 k0s.K0sConfig
	if err := yaml.Unmarshal(out, &kc2); err != nil {
		t.Fatalf("failed to re-unmarshal: %s", err.Error())
	}
	if kc2.DigString("spec", "api", "address") != "192.168.1.1" {
		t.Errorf("expected address 192.168.1.1, got %s", kc2.DigString("spec", "api", "address"))
	}
}
