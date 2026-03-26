# AI Agent Architecture Guide

This document provides **AI-specific guidance** for navigating and understanding the codebase. Use this to **prioritize files**, **analyze components**, and **propose changes** efficiently.

---

## 1. Key Files and Directories for AI Agents
AI agents should focus on the following files and directories to gather context quickly:

| Path | Purpose | AI-Specific Notes |
|------|---------|-------------------|
| [`pkg/engine/driver.go`](../../pkg/engine/driver.go) | Defines the `Driver` interface for implementations. | **Critical**: All providers (e.g., `k0s`, `MKE`) must implement this. |
| [`pkg/config/config.go`](../../pkg/config/config.go) | Defines the configuration schema. | **Critical**: Validate configs against this schema before proposing changes. |
| [`cmd/`](../../cmd/) | CLI entry points. | Start here to trace user commands (e.g., `launchpad apply`). |
| [`implementation/`](../../implementation/) | Provider-specific logic (e.g., `k0s`, `MKE`). | Focus on the `Driver` implementation for each provider. |
| [`pkg/engine/`](../../pkg/engine/) | Core orchestration logic. | **Critical**: Analyze this for changes to core workflows. |

---

## 2. AI-Specific Workflow Guidance
### 2.1 Analyzing the Codebase
1. **Start with Interfaces**:
   - Review `pkg/engine/driver.go` to understand provider contracts.
   - Review `pkg/config/config.go` to validate configurations.

2. **Trace CLI Commands**:
   - Begin in `cmd/` to identify entry points (e.g., `launchpad apply`).
   - Follow the flow to `pkg/engine/` and the relevant implementation (e.g., `implementation/k0s/`).

3. **Focus on Providers**:
   - Each provider in `implementation/` must adhere to the `Driver` interface.
   - Example: `implementation/k0s/driver.go` implements `Apply()`, `Destroy()`, etc.

### 2.2 Proposing Changes
- **Core Logic**: Propose changes to `pkg/engine/` for workflow modifications.
- **Providers**: Propose changes to `implementation/<provider>/` for provider-specific logic.
- **CLI**: Propose changes to `cmd/` for user-facing commands.

---

## 3. Critical Interfaces for AI Agents
AI agents must understand these interfaces to propose valid changes:

| Interface | Location | Methods | AI-Specific Notes |
|-----------|----------|---------|-------------------|
| `Driver` | `pkg/engine/driver.go` | `Apply()`, `Destroy()`, `Validate()` | All providers must implement these. |
| `Config` | `pkg/config/config.go` | `Validate()` | Validate configs against this schema. |

---

## 4. Low-Priority Tasks
The following tasks are **AI-driven enhancements** and should be treated as **low priority**:

- **AI-Generated TODOs**: Automatically adding TODOs or reminders based on code analysis.
- **AI-Assisted Refactoring**: Using AI to suggest or apply large-scale refactors.
- **Automated Documentation**: Generating or updating documentation without human review.

These features may be revisited in the future but are not a focus for current development.

---

## Next Steps
1. Review [`workflows.md`](./workflows.md) for common tasks (e.g., adding a feature).
2. Check [`code_guidelines.md`](./code_guidelines.md) for coding standards.