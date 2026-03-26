# AI Agent Workflows

This document outlines **AI-specific workflows** for analyzing, debugging, and proposing changes in this repository. Use this to **prioritize tasks** and **automate contributions** effectively.

---

## 1. Analyzing the Codebase
### Workflow: Gathering Context
1. **Start with Critical Interfaces**:
   - Review `pkg/engine/driver.go` to understand provider contracts.
   - Review `pkg/config/config.go` to validate configurations.

2. **Trace CLI Commands**:
   - Begin in `cmd/` to identify entry points (e.g., `launchpad apply`).
   - Follow the flow to `pkg/engine/` and the relevant implementation (e.g., `implementation/k0s/`).

3. **Focus on Providers**:
   - Each provider in `implementation/` must implement the `Driver` interface.
   - Example: `implementation/k0s/driver.go` implements `Apply()`, `Destroy()`, etc.

---

## 2. Debugging Issues
### Workflow: Assisting with Debugging
1. **Reproduce the Issue**:
   - Use the CLI to replicate the problem (e.g., `launchpad apply -f config.yaml`).
   - Check logs for errors (e.g., `launchpad --debug apply -f config.yaml`).

2. **Trace the Flow**:
   - Start in `cmd/` to identify the entry point.
   - Follow the flow to `pkg/engine/` and the relevant implementation (e.g., `implementation/k0s/`).

3. **Validate Configurations**:
   - Use `pkg/config/config.go` to validate YAML files.
   - Ensure the config matches the expected schema.

4. **Propose Fixes**:
   - Suggest changes to the relevant package (e.g., `pkg/engine/`, `implementation/k0s/`).
   - Ensure fixes align with the `Driver` interface and `Config` schema.

---

## 3. Proposing Changes
### Workflow: Assisting with Code Changes
1. **Gather Context**:
   - Read the relevant files in `docs/ai/` (e.g., `architecture.md`, `workflows.md`).
   - Review the codebase structure and key interfaces.

2. **Propose Changes**:
   - Use `architecture.md` to identify where changes should be made (e.g., `pkg/engine/`, `implementation/`).
   - Ensure proposals align with the `Driver` interface and `Config` schema.

3. **Validate Proposals**:
   - Run `go test ./...` to ensure no regressions.
   - Use `make lint` to check for style issues.

4. **Update Documentation**:
   - Suggest updates to `docs/ai/` if the codebase evolves (e.g., new providers, workflows).

---

## 4. Adding a New Provider
### Workflow: Assisting with Provider Implementation
1. **Create a New Directory**:
   - Add a subdirectory in `implementation/` (e.g., `implementation/new-provider/`).

2. **Implement the `Driver` Interface**:
   - Define methods like `Apply()`, `Destroy()`, and `Validate()` in `driver.go`.
   - Ensure the implementation adheres to `pkg/engine/driver.go`.

3. **Extend Configuration Support**:
   - Update `pkg/config/config.go` to include the new provider's schema.

4. **Register the Provider**:
   - Propose changes to the `engine` package to recognize the new provider.

5. **Add Tests**:
   - Suggest unit tests for the provider's package.
   - Propose integration tests in `integration/`.

---

## 5. AI-Specific Tasks
### Workflow: Automating Contributions
1. **Prioritize Critical Files**:
   - Focus on `pkg/engine/driver.go`, `pkg/config/config.go`, and `cmd/`.

2. **Propose Small, Incremental Changes**:
   - Avoid large-scale refactors unless explicitly requested.
   - Ensure changes are backward-compatible where possible.

3. **Update AI Documentation**:
   - Suggest updates to `docs/ai/` to reflect changes in the codebase.
   - Ensure `README.md` in `docs/ai/` remains the entry point for AI agents.