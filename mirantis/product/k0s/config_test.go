package k0s_test

import (
	"testing"

	"github.com/Mirantis/launchpad/mirantis/product/k0s"
	"gopkg.in/yaml.v3"
)

func Test_Config_K0SConfig(t *testing.T) {
	ctx := t.Context()
	ky := `
version: v1.28.4+k0s.0
config:
  apiVersion: k0s.k0sproject.io/v1beta1
  kind: ClusterConfig
  metadata:
    name: k0s
  spec:
    controllerManager: {}
    extensions:
      helm:
        concurrencyLevel: 5
        charts: null
        repositories: null
      storage:
        create_default_storage_class: true
        type: openebs_local_storage
    installConfig:
      users:
        etcdUser: etcd
        kineUser: kube-apiserver
        konnectivityUser: konnectivity-server
        kubeAPIserverUser: kube-apiserver
        kubeSchedulerUser: kube-scheduler
`

	var kn yaml.Node

	if err := yaml.Unmarshal([]byte(ky), &kn); err != nil {
		t.Errorf("k0s test yaml unmarshal fail: %s", err.Error())
	}

	kc, derr := k0s.DecodeComponent("test", kn.Decode)
	if derr != nil {
		t.Errorf("K0s decode failed: %s", derr.Error())
	}

	if verr := kc.Validate(ctx); verr != nil {
		t.Errorf("K0s component validate failed: %s", verr.Error())
	}
}
