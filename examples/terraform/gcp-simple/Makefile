.PHONY: apply clean launchpad-yaml

PROJECT_ID ?= $(shell gcloud config get-value project 2>/dev/null)
ACCESS_TOKEN ?= $(shell gcloud auth print-access-token 2>/dev/null)

apply:
	@if [ -z "$(PROJECT_ID)" ]; then echo "Error: PROJECT_ID not set and could not be detected via gcloud"; exit 1; fi
	@if [ -n "$(ACCESS_TOKEN)" ]; then \
		GOOGLE_OAUTH_ACCESS_TOKEN="$(ACCESS_TOKEN)" terraform apply -var="project_id=$(PROJECT_ID)" -auto-approve; \
	else \
		terraform apply -var="project_id=$(PROJECT_ID)" -auto-approve; \
	fi

launchpad-yaml:
	terraform output -raw launchpad_yaml > launchpad.yaml

up: apply launchpad-yaml

# Connect via generated SSH key
ssh:
	ssh -i id_rsa $(shell terraform output -raw ssh_user)@$(shell terraform output -raw public_ip)

# Connect via native gcloud (OS Login + IAP)
gcloud-ssh:
	gcloud compute ssh $(shell terraform output -raw instance_name) --zone=$(shell terraform output -raw zone) --tunnel-through-iap

# Apply launchpad project
launchpad-apply: launchpad-yaml
	../../../dist/launchpad project:apply --config launchpad.yaml

clean:
	@if [ -z "$(PROJECT_ID)" ]; then echo "Error: PROJECT_ID not set and could not be detected via gcloud"; exit 1; fi
	@if [ -n "$(ACCESS_TOKEN)" ]; then \
		GOOGLE_OAUTH_ACCESS_TOKEN="$(ACCESS_TOKEN)" terraform destroy -var="project_id=$(PROJECT_ID)" -auto-approve; \
	else \
		terraform destroy -var="project_id=$(PROJECT_ID)" -auto-approve; \
	fi
	rm -f launchpad.yaml id_rsa
