// ------- Ye old launchpad yaml (just for debugging)

locals {
  launchpad_yaml = <<-EOT
apiVersion: launchpad.mirantis.com/v2.1
kind: project
metadata:
  name: ${var.name}
spec:
  project:
    prune: false
  components:
    hosts:
%{~for h in local.k0s_hosts_ssh}
      # ${h.label} (ssh)
      ${h.id}:
        k0s:
          role: ${h.role}
        rig:
          ssh:
            address: ${h.ssh_address}
            user: ${h.ssh_user}
            keyPath: ${h.ssh_key_path}
%{~endfor}

    k0s:
      version: ${var.k0s.version}
      config:
        apiVersion: k0s.k0sproject.io/v1beta1
        kind: ClusterConfig
        metadata:
          creationTimestamp: null
          name: k0s
        spec:
          api:
            externalAddress: ${local.K0S_URL}
            k0sApiPort: 9443
            port: 6443
          extensions:
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
            type: etcd
          telemetry:
            enabled: true

    mke4:
      version: 4.0.0

    msr4:
      version: 4.0.0
EOT

}

output "launchpad_yaml" {
  description = "launchpad config file yaml (for debugging)"
  sensitive   = true
  value       = local.launchpad_yaml
}
