
LOCAL_TAG?=0.0.0

GOLANGCILINT?=$(shell which golangci-lint)
ifeq (, $(GOLANGCILINT))
	GOLANGCILINT=?docker run -ti --rm -v "$(CURDIR):/data" -w "/data" golangci/golangci-lint:latest golangci-lint
endif

# Lint by running golangci-lint in a docker container
.PHONY: lint
lint:
	$(GOLANGCILINT) run ./...

# Local install of the plugin
# @SEE README.md on how to use the locally built plugin
.PHONY: local
local:
	GORELEASER_CURRENT_TAG="$(LOCAL_TAG)" goreleaser build --clean --single-target --skip=validate --snapshot

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
