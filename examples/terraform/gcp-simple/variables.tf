variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP Region"
  type        = string
  default     = "us-east1"
}

variable "zone" {
  description = "GCP Zone"
  type        = string
  default     = "us-east1-b"
}

variable "name" {
  description = "Base name for resources"
  type        = string
  default     = "lp2-test-gcp"
}

variable "machine_type" {
  description = "Machine type for the VM"
  type        = string
  default     = "e2-standard-4" # 4 vCPUs, 16 GB RAM - suitable for MKE/MSR
}

variable "boot_disk_type" {
  description = "Type of the boot disk"
  type        = string
  default     = "pd-balanced"
}

variable "boot_disk_size" {
  description = "Size of the boot disk in GB"
  type        = number
  default     = 50
}

variable "spot" {
  description = "Whether to use a spot instance"
  type        = bool
  default     = true
}

variable "ssh_user" {
  description = "SSH user"
  type        = string
  default     = "ubuntu"
}

variable "additional_ssh_keys" {
  description = "Additional SSH keys to add to the VM (format: 'user:key')"
  type        = list(string)
  default     = []
}
