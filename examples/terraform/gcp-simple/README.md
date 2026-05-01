# GCP Simple Terraform Example

This Terraform configuration provisions a single GCP instance (Ubuntu 22.04) and generates a `launchpad.yaml` file for testing Launchpad.

## Usage

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Apply the configuration (provide your GCP Project ID):
   ```bash
   terraform apply -var="project_id=YOUR_PROJECT_ID"
   ```

   *Tip: Use `-var="spot=true"` (default) for cheaper testing instances.*

3. Extract the generated `launchpad.yaml`:
   ```bash
   terraform output -raw launchpad_yaml > launchpad.yaml
   ```

4. Run Launchpad:
   ```bash
   launchpad project apply
   ```

## Faster Testing Cycle

This stack is optimized for speed:
- **Spot Instances**: Enabled by default (`spot = true`) to keep costs low.
- **SSH Readiness**: Terraform waits for SSH to be actually available before finishing.
- **Cloud-Init**: Pre-installs common dependencies (`socat`, `nfs-common`, etc.) and disables `ufw` to avoid networking issues.
- **Balanced Disk**: Uses `pd-balanced` by default for better I/O performance than standard disks.

## Native GCP Access

While Launchpad uses the generated `id_rsa` key, you can also use native GCP tools for manual inspection (no local key management required):

```bash
# Connect via IAP tunnel (works even without public IP access)
gcloud compute ssh lp2-test-gcp --tunnel-through-iap
```

Or use the provided Makefile:
```bash
make gcloud-ssh
```

## Resources Created

- VPC and Subnetwork
- Firewall rules for SSH and Mirantis products
- Generated SSH keypair (saved as `id_rsa` in the current directory)
- Google Compute Instance (Ubuntu 22.04, e2-standard-4)
