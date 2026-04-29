package k0s_test

import (
	"testing"

	"github.com/k0sproject/version"

	"github.com/Mirantis/launchpad/mirantis/product/k0s"
	"gopkg.in/yaml.v3"
)

var (
	cfgy = `
  apiVersion: k0s.k0sproject.io/v1beta1
  kind: ClusterConfig
  metadata:
    creationTimestamp: null
    name: k0s
  spec:
    api:
      k0sApiPort: 9443
      port: 6443
      sans:
      - example-k0s-lb-c91e5f2cde668afd.elb.us-east-1.amazonaws.com
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
    konnectivity:
      adminPort: 8133
      agentPort: 8132
    network:
      calico: null
      clusterDomain: cluster.local
      dualStack: {}
      kubeProxy:
        iptables:
          minSyncPeriod: 0s
          syncPeriod: 0s
        ipvs:
          minSyncPeriod: 0s
          syncPeriod: 0s
          tcpFinTimeout: 0s
          tcpTimeout: 0s
          udpTimeout: 0s
        metricsBindAddress: 0.0.0.0:10249
        mode: iptables
      kuberouter:
        autoMTU: true
        hairpin: Enabled
        ipMasq: false
        metricsPort: 8080
        mtu: 0
        peerRouterASNs: ""
        peerRouterIPs: ""
      nodeLocalLoadBalancing:
        envoyProxy:
          apiServerBindPort: 7443
          konnectivityServerBindPort: 7132
        type: EnvoyProxy
      provider: kuberouter
    scheduler: {}
    storage:
      etcd:
        externalCluster: null
      type: etcd
    telemetry:
      enabled: true
`
)

// Test_ConfigURL verifies the download URL format without making network calls.
func Test_ConfigURL(t *testing.T) {
	arch := "amd64"
	vs := "v1.30.0+k0s.0"
	v, verr := version.NewVersion(vs)
	if verr != nil {
		t.Errorf("err with version %s: %s", vs, verr.Error())
	}

	url := k0s.DownloadK0sURL(*v, arch)
	expected := "https://github.com/k0sproject/k0s/releases/download/v1.30.0+k0s.0/k0s-v1.30.0+k0s.0-amd64"
	if url != expected {
		t.Errorf("unexpected URL: got %s, want %s", url, expected)
	}
}

func Test_K0SConfigModify(t *testing.T) {
	var cfg k0s.K0sConfig
	if err := yaml.Unmarshal([]byte(cfgy), &cfg); err != nil {
		t.Fatalf("fail to unmarshal test K0sConfig: %s", err.Error())
	}

	t.Logf("Initial CFG: %+v", cfg)

	if cfg.Dig("spec", "api", "k0sApiPort") == nil {
		t.Error("expected spec.api.k0sApiPort to be present")
	}
}
