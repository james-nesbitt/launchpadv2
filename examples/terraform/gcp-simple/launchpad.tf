locals {
  launchpad_yaml = <<-EOT
apiVersion: launchpad.mirantis.com/v2.1
kind: project
metadata:
  name: ${var.name}
spec:
  components:
    hosts:
      ${var.name}:
        rig:
          ssh:
            address: ${google_compute_instance.vm_instance.network_interface[0].access_config[0].nat_ip}
            user: ${var.ssh_user}
            keyPath: ${abspath(local_sensitive_file.private_key.filename)}
    k0s:
      version: 1.30.2+k0s.0
EOT
}
