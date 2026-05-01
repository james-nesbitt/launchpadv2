provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

resource "google_compute_network" "vpc_network" {
  name                    = "${var.name}-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "vpc_subnetwork" {
  name          = "${var.name}-subnet"
  ip_cidr_range = "10.0.1.0/24"
  region        = var.region
  network       = google_compute_network.vpc_network.id
}

resource "google_compute_firewall" "allow_ssh" {
  name    = "${var.name}-allow-ssh"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "allow_internal" {
  name    = "${var.name}-allow-internal"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  source_ranges = ["10.0.1.0/24"]
}

resource "google_compute_firewall" "allow_products" {
  name    = "${var.name}-allow-products"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["80", "443", "6443", "2376", "2377", "7946", "4789", "10250", "12376", "12379", "12380", "12381"]
  }

  allow {
    protocol = "udp"
    ports    = ["4789", "7946"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "allow_iap_ssh" {
  name    = "${var.name}-allow-iap-ssh"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["35.235.240.0/20"] # Official GCP IAP range for TCP forwarding
}

resource "tls_private_key" "ssh_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "local_sensitive_file" "private_key" {
  content         = tls_private_key.ssh_key.private_key_pem
  filename        = "${path.module}/id_rsa"
  file_permission = "0600"
}

resource "google_compute_instance" "vm_instance" {
  name         = var.name
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts"
      size  = var.boot_disk_size
      type  = var.boot_disk_type
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.vpc_subnetwork.id

    access_config {
      // Ephemeral public IP
    }
  }

  scheduling {
    preemptible                 = var.spot
    automatic_restart           = !var.spot
    provisioning_model          = var.spot ? "SPOT" : "STANDARD"
    instance_termination_action = var.spot ? "STOP" : null
  }

  metadata = {
    ssh-keys = join("\n", concat(
      ["${var.ssh_user}:${tls_private_key.ssh_key.public_key_openssh}"],
      var.additional_ssh_keys
    ))
    user-data = <<-EOT
      #cloud-config
      package_update: true
      package_upgrade: false
      packages:
        - apt-transport-https
        - ca-certificates
        - curl
        - gnupg
        - lsb-release
        - socat
        - nfs-common
      runcmd:
        - [ systemctl, stop, ufw ]
        - [ systemctl, disable, ufw ]
    EOT
  }
}

resource "terraform_data" "wait_for_ssh" {
  triggers_replace = [
    google_compute_instance.vm_instance.id
  ]

  connection {
    type        = "ssh"
    user        = var.ssh_user
    host        = google_compute_instance.vm_instance.network_interface[0].access_config[0].nat_ip
    private_key = tls_private_key.ssh_key.private_key_pem
  }

  provisioner "remote-exec" {
    inline = [
      "echo 'SSH is ready!'"
    ]
  }
}
