package k0s

import "github.com/k0sproject/dig"

/**
 * A Note About K0sConfig
 *
 * This should all be just a k0sProject/dig but in original development I had problems
 * building overrides after unmarshalling into a dig.  Overriding parts would drop
 * complete other parts of the config.
 *
 * We should return to using a dig, if we can get it to work.
 *
 */

type K0sConfig struct {
	APIVersion string            `yaml:"apiVersion" default:"k0s.k0sproject.io/v1beta1"`
	Kind       string            `yaml:"kind" default:"ClusterConfig"`
	Metadata   K0sConfigMetadata `yaml:"metadata"`
	Spec       K0sConfigSpec     `yaml:"spec"`
}

type K0sConfigMetadata struct {
	Name string `yaml:"name"`
}

type K0sConfigSpec struct {
	API     K0sConfigSpecApi     `yaml:"api"`
	Storage K0sConfigSpecStorage `yaml:"storage"`

	ControllerManager dig.Mapping                `yaml:"controllerManager"`
	Extensions        K0sConfigSpecExtensions    `yaml:"extensions,omitempty"`
	InstallConfig     K0sConfigSpecInstallConfig `yaml:"installConfig,omitempty"`
	//	Konnectivity      K0sConfigSpecKonnectivity  `yaml:"konnectivity,omitempty"`
	Network   dig.Mapping            `yaml:"network,omitempty"`
	Scheduler dig.Mapping            `yaml:"scheduler,omitempty"`
	Telemetry K0sConfigSpecTelemetry `yaml:"telemetry,omitempty"`
}

type K0sConfigSpecApi struct {
	ExternalAddress string   `yaml:"externalAddress"`
	Address         string   `yaml:"address"`
	K0sApiPort      int      `yaml:"k0sApiPort" default:"9443"`
	Port            int      `yaml:"port" default:"6443"`
	Sans            []string `yaml:"sans"`
}

type K0sConfigSpecStorage struct {
	Type string                   `yaml:"type"`
	Etcd K0sConfigSpecStorageEtcd `yaml:"etcd"`
}

type K0sConfigSpecStorageEtcd struct {
	PeerAddress     string      `yaml:"peerAddress"`
	externalCluster dig.Mapping `yaml:"externalCluster"`
}

type K0sConfigSpecExtensions struct {
	Helm    K0sConfigSpecExtensionsHelm   `yaml:"helm"`
	Storage K0sConfigSpecExtensionStorage `yaml:"storage"`
}

type K0sConfigSpecExtensionsHelm struct {
	ConcurrencyLevel int                                         `yaml:"concurrencyLevel"`
	Repositories     map[string]string                           `yaml:"repositories,omitempty"`
	Charts           map[string]K0sConfigSpecExtensionsHelmChart `yaml:"charts,omitempty"`
}

type K0sConfigSpecExtensionStorage struct {
	Type                      string `yaml:"type"`
	CreateDefaultStorageClass bool   `yaml:"create_default_storage_class"`
}

type K0sConfigSpecExtensionsHelmChart struct {
}

type K0sConfigSpecInstallConfig struct {
	Users map[string]string `yaml:"users"`
}

type K0sConfigSpecKonnectivity struct {
	AdminPort int `yaml:"adminPort,omitempty"`
	AgentPort int `yaml:"agentPort,omitempty"`
}

type K0sConfigSpecTelemetry struct {
	Enabled bool `yaml:"enabled"`
}

// # https://docs.k0sproject.io/head/configuration/#using-a-configuration-file
//apiVersion: k0s.k0sproject.io/v1beta1
//kind: ClusterConfig
//metadata:
//  name: k0s
//spec:
//  api:
//    address: 192.168.68.104
//    k0sApiPort: 9443
//    port: 6443
//    sans:
//    - 192.168.68.104
//  controllerManager: {}
//  extensions:
//    helm:
//      concurrencyLevel: 5
//  installConfig:
//    users:
//      etcdUser: etcd
//      kineUser: kube-apiserver
//      konnectivityUser: konnectivity-server
//      kubeAPIserverUser: kube-apiserver
//      kubeSchedulerUser: kube-scheduler
//  konnectivity:
//    adminPort: 8133
//    agentPort: 8132
//  network:
//    clusterDomain: cluster.local
//    dualStack:
//      enabled: false
//    kubeProxy:
//      iptables:
//        minSyncPeriod: 0s
//        syncPeriod: 0s
//      ipvs:
//        minSyncPeriod: 0s
//        syncPeriod: 0s
//        tcpFinTimeout: 0s
//        tcpTimeout: 0s
//        udpTimeout: 0s
//      metricsBindAddress: 0.0.0.0:10249
//      mode: iptables
//    kuberouter:
//      autoMTU: true
//      hairpin: Enabled
//      metricsPort: 8080
//    nodeLocalLoadBalancing:
//      enabled: false
//      envoyProxy:
//        apiServerBindPort: 7443
//        konnectivityServerBindPort: 7132
//      type: EnvoyProxy
//    podCIDR: 10.244.0.0/16
//    provider: kuberouter
//    serviceCIDR: 10.96.0.0/12
//  scheduler: {}
//  storage:
//    etcd:
//      peerAddress: 192.168.68.104
//    type: etcd
//telemetry:
//    enabled: true
