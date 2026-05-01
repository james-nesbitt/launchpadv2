
LOCAL_TAG?=0.0.0

# Lint by running golangci-lint
.PHONY: lint
lint:
	golangci-lint run ./...

# Local build of the binary
.PHONY: local
local:
	go build -o dist/launchpad ./mirantis/cmd/launchpad/main.go

dist:
	GORELEASER_CURRENT_TAG="$(LOCAL_TAG)" goreleaser build --clean --skip=validate

.PHONY: clean
clean:
	rm -rf dist

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...

# Run govulncheck to scan for vulnerabilities
.PHONY: security
security:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...
