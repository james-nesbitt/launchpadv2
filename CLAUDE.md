# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make test        # Run all unit tests
make lint        # Run golangci-lint (uses Docker if not installed locally)
make fmt         # Format all Go files
make local       # Build binary for current platform via GoReleaser
make dist        # Build binaries for all platforms
make security    # Run govulncheck for vulnerability scanning
make clean       # Remove dist directory
```

**Run a single test:**
```bash
go test -v -run TestName ./path/to/package
# Example: go test -v -run Test_CommandBuild ./pkg/action
```

**Run integration tests** (require `//go:build integration` tag):
```bash
go test -tags integration ./...
```

## Architecture

Launchpad is a CLI tool for installing and managing Mirantis cluster products (MKE, MSR, MCR, K0s). It uses a **component-based, interface-driven architecture** where each product is a modular component that participates in command orchestration via a dependency graph.

### Execution Flow

```
CLI Command → Bootstrap (cmd/bootstrap.go)
           → Load YAML config (mirantis/config/v2_0/ or v2_1/)
           → Decode components (pkg/component/, mirantis/product/)
           → Resolve dependencies (pkg/dependency/)
           → Build command phases (pkg/action/)
           → Execute phases sequentially
```

The main entry point is `mirantis/cmd/launchpad/main.go`, which registers all products via blank `init()` imports.

### Core Concepts

**Component** — A product (e.g., MKE4, K0s) that:
- Implements `action.CommandBuild` to add phases to commands (`apply`, `reset`, `discover`)
- Optionally implements `dependency.RequiresDependencies` / `dependency.ProvidesDependencies`

**Implementation** — Shared infrastructure (e.g., Kubernetes API, Docker, Rig host) used by multiple components. Lives in `implementation/`.

**Dependency System** — Components declare abstract requirements (e.g., "Kubernetes API") that other components fulfill. The engine resolves these at runtime via a dependency graph.

**Config versions** — Two YAML spec versions are supported: `v2.0` (`mirantis/config/v2_0/`) and `v2.1` (`mirantis/config/v2_1/`). The v2.1 spec is current.

### Key Package Responsibilities

| Package | Role |
|---------|------|
| `pkg/action/` | Command and phase orchestration |
| `pkg/component/` | Component registry and decoding |
| `pkg/dependency/` | Requirement/fulfillment resolution |
| `pkg/project/` | Project management, CLI building |
| `pkg/host/` | Host abstraction layer |
| `mirantis/product/` | Mirantis product implementations (k0s, mke3, mke4, msr2, msr3, msr4, mcr) |
| `implementation/` | Shared implementations (docker, kubernetes, rig, hook) |
| `mirantis/config/` | YAML configuration parsing for v2.0 and v2.1 |

### Adding a New Component

1. Define a `ProductDecoder` and register it in `init()`:
   ```go
   func init() {
       product.RegisterDecoder("my-component", NewMyComponentDecoder)
   }
   ```
2. Implement `action.CommandBuild` to contribute phases to commands.
3. Optionally implement `dependency.RequiresDependencies` / `dependency.ProvidesDependencies`.
4. Import the new package with a blank import in `mirantis/cmd/launchpad/main.go`.

### Logging

Use `log/slog` (stdlib structured logging) — not `log.Printf` or other logging packages.

### Error Handling

Wrap errors with context using `fmt.Errorf("failed to do X: %w", err)`.

## Extended Documentation

- `docs/design.md` — Core concepts and architecture diagrams
- `docs/component.md` — Component interface guide with examples
- `docs/config.md` — YAML configuration schema
- `docs/commands.md` — CLI command reference
- `docs/ai/` — AI-agent-specific guidance (architecture, workflows, code guidelines)
