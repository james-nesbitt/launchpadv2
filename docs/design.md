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

**Example**:
```mermaid
graph LR
  A[Component A] -->|Phase 1| B[Component B]
  B -->|Phase 2| C[Component C]
```

---

### 4. Dependency System
**Requirement**: A need (e.g., "Kubernetes cluster").
**Dependency**: A fulfillment (e.g., "MKE provides Kubernetes").

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

**Example**:
```go
project := &project.Project{
  Components: []component.Component{
    mke.NewComponent(),
    msr.NewComponent(),
  },
}
```