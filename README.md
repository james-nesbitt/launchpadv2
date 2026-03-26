# Mirantis Launchpad

**A modular tool for installing and managing Mirantis cluster products (MKE, MSR, K0s) with flexibility and extensibility at its core.**

Launchpad is designed to:
1. **Facilitate strategic and tactical management** of Mirantis products.
2. **Stay flexible** as products and customer needs evolve.
3. **Remain extensible** through a component-based architecture.

---

## Core Purpose
Launchpad enables **modular, interface-driven management** of Mirantis clusters by:
- **Decoupling components** via registered capabilities and Go interfaces.
- **Swapping implementations** at compile time (e.g., MKE vs. K0s for Kubernetes).
- **Automating dependencies** (e.g., host roles, APIs) through abstract requirements.

This ensures **adaptability** to product changes and **customizability** for diverse use cases.

---

## Key Features
### 1. **Modular Architecture**
- **Components**: Functional units (e.g., MKE, MSR) that implement interfaces like `action.CommandBuild`.
- **Dependencies**: Abstract requirements (e.g., "Kubernetes API") fulfilled by components.
- **Extensibility**: Add new products or features by registering components.

### 2. **Strategic & Tactical Management**
- **Strategic**: Define cluster configurations in YAML (e.g., `cluster.yaml`).
- **Tactical**: Execute commands like `apply`, `discover`, and `reset`.

### 3. **Flexibility**
- **Compile-Time Swapping**: Replace components (e.g., MKE → K0s) without rewriting core logic.
- **Customer-Centric**: Adapt to unique environments (e.g., air-gapped, multi-cloud).

---

## Getting Started
### 1. **Installation**
#### Build from Source
```bash
# Build for all platforms (requires GoReleaser)
make clean dist

# Build for your local platform
make local
```

#### Download Binaries
Pre-built binaries are available in [GitHub Releases](https://github.com/Mirantis/launchpad/releases).

### 2. **Quick Start**
1. **Define a cluster** in `cluster.yaml`:
   ```yaml
   cluster:
     name: "my-cluster"
   components:
     - name: "mke"
       version: "3.6.0"
   ```
2. **Apply the configuration**:
   ```bash
   launchpad apply --config cluster.yaml
   ```

---

## Documentation
### Core
- [Commands](./docs/commands.md): CLI usage and examples.
- [Configuration](./docs/config.md): YAML schema and validation.
- [Components](./docs/component.md): Extending Launchpad with modular components.
- [Design](./docs/design.md): Architecture, interfaces, and dependency system.

### Development
- [TODO](./docs/TODO.md): Backlog and roadmap.
- [Examples](./examples/): Sample configurations and scripts.

---

## Contributing
1. **Report Issues**: Open a GitHub issue for bugs or feature requests.
2. **Submit PRs**: Follow the [contribution guidelines](./CONTRIBUTING.md).
3. **Join Discussions**: Share feedback in [Mirantis Community Slack](https://mirantiscommunity.slack.com).

---

## License
Launchpad is licensed under the [Apache License 2.0](./LICENSE).