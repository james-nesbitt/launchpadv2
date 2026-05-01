output "public_ip" {
  description = "The public IP address of the instance"
  value       = google_compute_instance.vm_instance.network_interface[0].access_config[0].nat_ip
}

output "ssh_user" {
  description = "The SSH user"
  value       = var.ssh_user
}

output "instance_name" {
  description = "The name of the instance"
  value       = google_compute_instance.vm_instance.name
}

output "zone" {
  description = "The zone of the instance"
  value       = google_compute_instance.vm_instance.zone
}

output "ssh_command" {
  description = "The SSH command to connect to the instance using the generated key"
  value       = "ssh -i id_rsa ${var.ssh_user}@${google_compute_instance.vm_instance.network_interface[0].access_config[0].nat_ip}"
}

output "gcloud_ssh_command" {
  description = "The gcloud command to connect to the instance using OS Login / IAP"
  value       = "gcloud compute ssh ${google_compute_instance.vm_instance.name} --zone=${google_compute_instance.vm_instance.zone} --tunnel-through-iap"
}

output "launchpad_yaml" {
  description = "Generated launchpad.yaml for this host"
  value       = local.launchpad_yaml
  sensitive   = true
}
