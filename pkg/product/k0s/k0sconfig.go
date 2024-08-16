package k0s

import "github.com/k0sproject/dig"

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

//    controllerManager: {}
//    extensionsb:
//      helm:
//        charts: null
//        concurrencyLevel: 5
//        repositories: null
//      storage:
//        create_default_storage_class: false
//        type: external_storage
//    installConfig:
//      users:
//        etcdUser: etcd
//        kineUser: kube-apiserver
//        konnectivityUser: konnectivity-server
//        kubeAPIserverUser: kube-apiserver
//        kubeSchedulerUser: kube-scheduler
//    konnectivity:
//      adminPort: 8133
//      agentPort: 8132
//    network:
//      calico: null
//      clusterDomain: cluster.local
//      dualStack: {}
//      kubeProxy:
//        iptables:
//          minSyncPeriod: 0s
//          syncPeriod: 0s
//        ipvs:
//          minSyncPeriod: 0s
//          syncPeriod: 0s
//          tcpFinTimeout: 0s
//          tcpTimeout: 0s
//          udpTimeout: 0s
//        metricsBindAddress: 0.0.0.0:10249
//        mode: iptables
//      kuberouter:
//        autoMTU: true
//        hairpin: Enabled
//        ipMasq: false
//        metricsPort: 8080
//        mtu: 0
//        peerRouterASNs: ""
//        peerRouterIPs: ""
//      nodeLocalLoadBalancing:
//        envoyProxy:
//          apiServerBindPort: 7443
//          konnectivityServerBindPort: 7132
//        type: EnvoyProxy
//      provider: kuberouter
//    scheduler: {}
//    storage:
//      etcd:
//        externalCluster: null
//      type: etcd
//    telemetry:
//      enabled: true
