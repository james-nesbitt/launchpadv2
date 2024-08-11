// constants
locals {

  // role for MSR machines, so that we can detect if msr config is needed
  msr_role = "msr"
  // only hosts with these roles will be used for launchpad_yaml
  mcr_roles = ["manager", "worker", local.msr_role]

}

variable "mcr" {
  description = "MCR configuration"
  type = object({
    version = string
  })
  default = {
    version = "23.0.9"
  }
}

variable "mke" {
  description = "MKE configuration"
  type = object({
    version = string
    connect = object({
      username = string
      password = string
      insecure = optional(bool, false) // true if this endpoint will not use a valid certificate      
    })
  })
  default = {
    version = "3.7.6"
    connect = {
      username = "admin"
      password = ""
    }
  }
}

variable "msr" {
  description = "MSR configuration"
  type = object({
    version = string
  })
  default = {
    version = ""
  }
}


// locals calculated before the provision run
locals {
  // standard MKE ingresses
  mke_ingresses = {
    "mke" = {
      description = "MKE ingress for UI and Kube"
      nodegroups  = [for k, ng in var.nodegroups : k if ng.role == "manager"]

      routes = {
        "mke" = {
          port_incoming = 443
          port_target   = 443
          protocol      = "TCP"
        }
        "kube" = {
          port_incoming = 6443
          port_target   = 6443
          protocol      = "TCP"
        }
      }
    }
  }

  // standard MCR/MKE/MSR firewall rules [here we just leave it open until we can figure this out]
  mkex_securitygroups = {
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

output "mke_connect" {
  description = "Connection information for connecting to MKE"
  sensitive   = true
  value = {
    host     = local.MKE_URL
    username = var.mke.connect.username
    password = var.mke.connect.password
    insecure = var.mke.connect.insecure
  }
}
