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

/**
 * K0sConfig is a flexible mapping for k0s configuration.
 *
 * It uses github.com/k0sproject/dig to allow full flexibility in
 * the configuration structure while maintaining a core known structure
 * for common fields like apiVersion and kind.
 *
 * Example YAML:
 *
 * apiVersion: k0s.k0sproject.io/v1beta1
 * kind: ClusterConfig
 * spec:
 *   api:
 *     port: 6443
 *   storage:
 *     type: etcd
 */
type K0sConfig struct {
	APIVersion string      `yaml:"apiVersion" default:"k0s.k0sproject.io/v1beta1"`
	Kind       string      `yaml:"kind" default:"ClusterConfig"`
	Metadata   dig.Mapping `yaml:"metadata,omitempty"`
	Spec       dig.Mapping `yaml:"spec,omitempty"`
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
