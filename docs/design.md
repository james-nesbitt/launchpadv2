# Launchpad Design

## Core Purpose
Launchpad is designed to:
1. **Facilitate management** of Mirantis products (MKE, MSR, K0s).
2. **Stay flexible** as products and customer needs evolve.
3. **Remain extensible** through a component-based architecture.

This is achieved through **interface-driven decoupling** and **abstract dependencies**.

---

## Core Concepts

### 1. Component
**Definition**: A modular unit that:
- Builds commands (e.g., `apply`).
- Provides/requires dependencies.
- Implements interfaces (e.g., `action.CommandBuild`).

**Example**: MKE, MSR, or custom components.

---

### 2. Implementation
**Definition**: Shared functionality (e.g., Kubernetes API) used by multiple components.
**Example**: Both MKE and K0s components may provide the same Kubernetes implementation.

---

### 3. Command
**Definition**: A single operation (e.g., `launchpad apply`) built from component phases.

**Flow**:
1. Launchpad collects phases from components.
2. Orders phases by dependencies.
3. Executes phases sequentially.

**Diagram: Command Execution Flow**
```mermaid
graph LR
  A[CLI Command] -->|Parse Config| B[Engine]
  B -->|Collect Phases| C[Component A]
  C -->|Phase 1| D[Component B]
  D -->|Phase 2| E[Component C]
  E -->|Result| A
```

---

### 4. Dependency System
**Requirement**: A need (e.g., "Kubernetes cluster").
**Dependency**: A fulfillment (e.g., "MKE provides Kubernetes").

**Diagram: Dependency Resolution**
```mermaid
graph TD
  A[Component A] -->|Requires| B[Kubernetes >=1.20]
  C[Component B] -->|Provides| B
  B -->|Resolved| D[Engine]
```

**Example**:
```go
// Component A requires Kubernetes.
req := kubernetes.NewRequirement(">=1.20")

// Component B provides Kubernetes.
dep := &KubernetesDependency{Version: "1.21"}
```

---

### 5. Project
**Definition**: A collection of components and configurations.

**Diagram: Project Structure**
```mermaid
classDiagram
  class Project {
    +Components[] Component
    +Config Config
  }
  class Component {
    +Name string
    +Phases[] Phase
  }
  Project --> Component
```

**Example**:
```go
project := &project.Project{
  Components: []component.Component{
    mke.NewComponent(),
    msr.NewComponent(),
  },
}
```

---

## Architecture Overview

### Diagram: High-Level Architecture
```mermaid
graph TD
  A[CLI] -->|Commands| B[Engine]
  B -->|Orchestrates| C[Components]
  C -->|Provides/Requires| D[Dependencies]
  B -->|Validates| E[Config]
  C -->|Implements| F[Implementations]
```

### Key Components:
| Component | Purpose | Location |
|-----------|---------|----------|
| CLI | User interface | `cmd/` |
| Engine | Core orchestration | `pkg/engine/` |
| Config | Configuration parsing | `pkg/config/` |
| Components | Modular units (MKE, MSR) | `implementation/` |
| Dependencies | Requirement resolution | `pkg/engine/dependencies/` |
