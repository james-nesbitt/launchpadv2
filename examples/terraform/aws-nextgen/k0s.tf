// constants
locals {

  // only hosts with these roles will be used for k0s_yaml
  k0s_roles = ["controller", "worker"]

}

// Launchpad configuration
variable "k0s" {
  description = "k0s install configuration"
  type = object({
    version = string
  })
}

// locals calculated before the provision run
locals {
  // standard K0s ingresses
  k0s_ingresses = {
    "k0s" = {
      description = "K0s ingress for Kube"
      nodegroups  = [for k, ng in var.nodegroups : k if ng.role == "controller"]

      routes = {
        "kube" = {
          port_incoming = 6443
          port_target   = 6443
          protocol      = "TCP"
        }
      }
    }
  }

  // standard k0s firewall rules [here we just leave it open until we can figure this out]
  k0s_securitygroups = {
    "permissive" = {
      description = "Common SG for all project machines"
      nodegroups  = [for n, ng in var.nodegroups : n]
      ingress_ipv4 = [
        {
          description : "Permissive internal traffic [BAD RULE]"
          from_port : 0
          to_port : 0
          protocol : "-1"
          self : true
          cidr_blocks : []
        },
        {
          description : "Permissive external traffic [BAD RULE]"
          from_port : 0
          to_port : 0
          protocol : "-1"
          self : false
          cidr_blocks : ["0.0.0.0/0"]
        }
      ]
      egress_ipv4 = [
        {
          description : "Permissive outgoing traffic"
          from_port : 0
          to_port : 0
          protocol : "-1"
          cidr_blocks : ["0.0.0.0/0"]
          self : false
        }
      ]
    }
  }

}

// prepare values to make it easier to feed into launchpad
locals {
  // The SAN URL for the MKE load balancer ingress that is for the MKE load balancer
  K0S_URL = module.provision.ingresses["k0s"].lb_dns

  // flatten nodegroups into a set of objects with the info needed for each node, by combining the group details with the node detains
  k0s_hosts_ssh = merge([for k, ng in local.nodegroups : { for l, ngn in ng.nodes : ngn.label => {
    id : "${k}_${l}"
    label : ngn.label
    role : ng.role

    address : ngn.public_address

    ssh_address : ngn.public_ip
    ssh_user : ng.ssh_user
    ssh_port : ng.ssh_port
    ssh_key_path : abspath(local_sensitive_file.common_private_key.filename)
  } if contains(local.k0s_roles, ng.role) && ng.connection == "ssh" }]...)

}
