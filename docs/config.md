# Configuration

## Overview
Launchpad uses YAML files to define cluster configurations.
**Example**: [`examples/basic.yaml`](./examples/basic.yaml)

---

## Schema
### Diagram: Configuration Hierarchy
```mermaid
classDiagram
  class Config {
    +cluster Cluster
    +components[] Component
    +hosts[] Host
    +ai AI
  }
  class Cluster {
    +name string
  }
  class Component {
    +name string
    +version string
  }
  class Host {
    +role string
    +address string
    +ssh SSH
  }
  class SSH {
    +user string
    +keyPath string
  }
  class AI {
    +optimize bool
    +debug bool
  }
  Config --> Cluster
  Config --> Component
  Config --> Host
  Host --> SSH
  Config --> AI
```

### Required Fields
```yaml
cluster:
  name: "my-cluster"
components:
  - name: "mke"
    version: "3.6.0"
  - name: "msr"
    version: "2.9.0"
```

### Advanced Fields
```yaml
hosts:
  - role: "manager"
    address: "192.168.1.10"
    ssh:
      user: "ubuntu"
      keyPath: "~/.ssh/id_rsa"
```

### AI-Optimized Fields (Optional)
```yaml
ai:
  optimize: true  # Enable AI-driven optimization during `apply`.
  debug: true     # Enable AI-driven troubleshooting during `discover`.
```

---

## Validation
**Command**:
```bash
launchpad validate --config <file.yaml>
```

**Diagram: Validation Flow**
```mermaid
flowchart TD
  A[Config File] -->|Parse| B[Schema Validation]
  B -->|Check| C[Component Compatibility]
  C -->|Check| D[Dependency Fulfillment]
  D -->|Check| E[SSH Connectivity]
  E -->|Optional| F[AI Optimization Suggestions]
  F -->|Result| G[Success/Failure]
```

**Checks**:
- Component compatibility.
- Dependency fulfillment.
- SSH connectivity.
- AI optimization suggestions (if enabled).

---

## Validation
**Command**:
```bash
launchpad validate --config <file.yaml>
```
**Checks**:
- Component compatibility.
- Dependency fulfillment.
- SSH connectivity.
- AI optimization suggestions (if enabled).