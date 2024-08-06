package k0s_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/k0sproject/version"

	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/Mirantis/launchpad/pkg/product/k0s"
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

func NotTest_ConfigURL(t *testing.T) {
	arch := "amd64"
	vs := "v1.30.0+k0s.0"
	v, verr := version.NewVersion(vs)
	if verr != verr {
		t.Errorf("err with version %s: %s", vs, verr.Error())
	}

	url := k0s.DownloadK0sURL(*v, arch)

	r, err := http.Get(url)
	if err != nil {
		t.Errorf("couldn't download version:")
	}
	r.Body.Close()
}

func NotTest_BuildHostConfig(t *testing.T) {
	ctx := context.Background()
	ny := `
network:
  private_address: 192.168.100.10
  public_address: 100.1.1.10
`
	ky := `
role: controller
`

	h := host.NewHost("dummy")

	var nc yaml.Node
	if err := yaml.Unmarshal([]byte(ny), &nc); err != nil {
		t.Fatalf("fail to unmarshal mock host plugin yaml: %s", err.Error())
	}
	if err := h.DecodeHostPlugin(ctx, mock.HostRoleMock, nc.Decode); err != nil {
		t.Fatalf("fail to decode k0s host plugin: %s", err.Error())
	}

	var kc yaml.Node
	if err := yaml.Unmarshal([]byte(ky), &kc); err != nil {
		t.Fatalf("fail to unmarshal test yaml: %s", err.Error())
	}
	if err := h.DecodeHostPlugin(ctx, k0s.ComponentType, kc.Decode); err != nil {
		t.Fatalf("fail to decode k0s host plugin: %s", err.Error())
	}

	var cfg k0s.K0sConfig
	if err := yaml.Unmarshal([]byte(cfgy), &cfg); err != nil {
		t.Fatalf("fail to decode k0s yaml: %s", err.Error())
	}

	sans := []string{
		"san.example.org",
	}

	kh := k0s.HostGetK0s(h)

	cfg, cfgerr := kh.BuildHostConfig(ctx, cfg, sans)
	if cfgerr != nil {
		t.Fatalf("fail to build host config: %s", cfgerr.Error())
	}
}

func Test_K0SConfigModify(t *testing.T) {
	var cfg k0s.K0sConfig
	if err := yaml.Unmarshal([]byte(cfgy), &cfg); err != nil {
		t.Fatalf("fail to unmarshal test K0sConfig: %s", err.Error())
	}

	t.Logf("Initial CFG: %+v", cfg)

}
