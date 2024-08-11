// prepare values to make it easier to feed into launchpad
locals {
  // The SAN URL for the MKE load balancer ingress that is for the MKE load balancer
  MKE_URL = module.provision.ingresses["mke"].lb_dns

  // flatten nodegroups into a set of objects with the info needed for each node, by combining the group details with the node detains
  launchpad_hosts_ssh = merge([for k, ng in local.nodegroups : { for l, ngn in ng.nodes : ngn.label => {
    id : "${k}_${l}"
    label : ngn.label
    role : ng.role

    address : ngn.public_address

    ssh_address : ngn.public_ip
    ssh_user : ng.ssh_user
    ssh_port : ng.ssh_port
    ssh_key_path : abspath(local_sensitive_file.common_private_key.filename)
  } if contains(local.launchpad_roles, ng.role) && ng.connection == "ssh" }]...)
  launchpad_hosts_winrm = merge([for k, ng in local.nodegroups : { for l, ngn in ng.nodes : ngn.label => {
    id : "${k}_${l}"
    label : ngn.label
    role : ng.role

    address : ngn.public_address

    winrm_address : ngn.public_ip
    winrm_user : ng.winrm_user
    winrm_password : var.windows_password
    winrm_useHTTPS : ng.winrm_useHTTPS
    winrm_insecure : ng.winrm_insecure
  } if contains(local.launchpad_roles, ng.role) && ng.connection == "winrm" }]...)

  // decide if we need msr configuration (the [0] is needed to prevent an error of no msr instances exit)
  has_msr = sum(concat([0], [for k, ng in local.nodegroups : ng.count if ng.role == local.launchpad_role_msr])) > 0
}


// ------- Ye old launchpad yaml (just for debugging)

locals {
  launchpad_yaml_21 = <<-EOT
apiVersion: launchpad.mirantis.com/v2.1
kind: project
metadata:
  name: ${var.name}
spec:
  project:
    prune: false

  components:
    hosts:
%{~for h in local.launchpad_hosts_ssh}
      # ${h.label} (ssh)
      ${h.id}:
        mcr:
          role: ${h.role}
        rig:
          ssh:
            address: ${h.ssh_address}
            user: ${h.ssh_user}
            keyPath: ${h.ssh_key_path}
%{~endfor}
%{~for h in local.launchpad_hosts_winrm}
      # ${h.label} (winrm)
      ${h.id}:
        mcr:
          role: ${h.role}
        rig:
          winRM:
            address: ${h.winrm_address}
            user: ${h.winrm_user}
            password: ${h.winrm_password}
            useHTTPS: ${h.winrm_useHTTPS}
            insecure: ${h.winrm_insecure}
%{~endfor}

    mcr:
      version: ${var.mcr.version}
      repoURL: https://repos.mirantis.com
      installURLLinux: https://get.mirantis.com/
      installURLWindows: https://get.mirantis.com/install.ps1
      channel: stable

    mke3:
      version: ${var.mke.version}
      imageRepo: docker.io/mirantis
      install:
        adminUsername: "${var.mke.connect.username}"
        adminPassword: "${var.mke.connect.password}"
        san: "${local.MKE_URL}"
        flags: 
        - "--default-node-orchestrator=kubernetes"
        - "--nodeport-range=32768-35535"
      upgrade:
        flags:
        - "--force-recent-backup"
        - "--force-minimums"
      prune: true
%{if local.has_msr}

    msr2:
      version: ${var.msr.version}
      imageRepo: docker.io/mirantis
      "replicaIDs": "sequential"
      installFlags:
      - "--ucp-insecure-tls"
%{endif}
EOT

}

output "launchpad_yaml" {
  description = "launchpad config file yaml (for debugging)"
  sensitive   = true
  value       = local.launchpad_yaml_20
}
