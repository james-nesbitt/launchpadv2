package k0s_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/product/k0s"
	"gopkg.in/yaml.v3"
)

func Test_K0sConfigUnMarshall(t *testing.T) {
	k0sConfigYaml := `
apiVersion: k0s.k0sproject.io/v1beta1
kind: ClusterConfig
metadata:
  creationTimestamp: null
  name: k0s
spec:
  api:
    externalAddress: my.lb.org
    k0sApiPort: 9443
    port: 6443
    sans:
    - my.san.org
  controllerManager: {}
  extensions:
    helm:
      charts: null
      concurrencyLevel: 5
      repositories: null
    storage:
      create_default_storage_class: false
      type: external_storage
  installConfig:
    users:
      etcdUser: etcd
      kineUser: kube-apiserver
      konnectivityUser: konnectivity-server
      kubeAPIserverUser: kube-apiserver
      kubeSchedulerUser: kube-scheduler
  storage:
    etcd:
      externalCluster: null
      peerAddress: 172.31.2.62
    type: etcd
  telemetry:
    enabled: true`

	kc := k0s.K0sConfig{}

	if err := yaml.Unmarshal([]byte(k0sConfigYaml), &kc); err != nil {
		t.Errorf("failed to unmarshal: %s", err.Error())
	}

	if kc.Spec.API.ExternalAddress != "my.lb.org" {
		t.Errorf("wrong spec.api.externalAddress: %s", kc.Spec.API.ExternalAddress)
	}
	if kc.Spec.API.Sans[0] != "my.san.org" {
		t.Errorf("wrong spec.api.sans[0]: %s", kc.Spec.API.Sans)
	}
}
