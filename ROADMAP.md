# LaunchpadV2 Roadmap

This roadmap outlines the phased approach for evolving LaunchpadV2 to support current and future Mirantis products.

---

## Phase 1: Architecture and Design Review & Improvements
**Goal**: Ensure the codebase is scalable, maintainable, and aligned with best practices for supporting current and future Mirantis products.

### Key Deliverables:
1. **Codebase Audit**
   - Review `pkg/`, `implementation/`, and `cmd/` for technical debt, inconsistencies, and gaps.
   - Document findings in `docs/architecture-audit.md`.

2. **Design Improvements**
   - Refactor for modularity (e.g., pluggable providers for Mirantis products).
   - Improve error handling, logging, and configuration management.
   - Add unit/integration test coverage for critical paths.

3. **Architecture Documentation**
   - Update `docs/architecture.md` to reflect current and future state.
   - Include diagrams for key workflows (e.g., cluster lifecycle, product integration).

4. **Tooling and Automation**
   - Enhance `Makefile` for common tasks (e.g., linting, testing, building).
   - Add GitHub Actions or GitLab CI for automated testing and releases.

---

## Phase 2: Existing Implementation Wrap-Up
**Goal**: Finalize support for existing Mirantis products and ensure stability.

### Supported Products:
- **MCR** (Mirantis Container Runtime)
- **MKE3** (Mirantis Kubernetes Engine v3)
- **MSR2** (Mirantis Secure Registry v2)
- **MSR3** (Mirantis Secure Registry v3)
- **MSR4** (Mirantis Secure Registry v4)
- **k0s** (Zero-friction Kubernetes)

### Key Deliverables:
1. **Feature Completeness**
   - Validate all CLI commands and workflows for existing products.
   - Fix critical bugs and edge cases.

2. **Testing**
   - Expand integration tests in `integration/` for all supported products.
   - Add end-to-end tests for common deployment scenarios.

3. **Documentation**
   - Update `README.md` and `docs/` with usage examples for each product.
   - Add troubleshooting guides for common issues.

4. **Release Preparation**
   - Tag a stable release (e.g., `v1.0.0`) for existing product support.
   - Publish release notes and changelog.

---

## Phase 3: New Mirantis Product Implementation
**Goal**: Extend LaunchpadV2 to support new Mirantis products.

### Target Products:
- **MKE4** (Mirantis Kubernetes Engine v4)
- **k0rdent** (Enterprise-grade k0s distribution)

### Key Deliverables:
1. **Requirements Gathering**
   - Collaborate with product teams to define integration requirements.
   - Document use cases and workflows in `docs/new-products.md`.

2. **Implementation**
   - Add new providers in `implementation/` for MKE4 and k0rdent.
   - Extend CLI commands and APIs to support new features.

3. **Testing**
   - Develop unit and integration tests for new products.
   - Validate interoperability with existing products (e.g., MKE4 + MSR4).

4. **Documentation**
   - Add usage guides for new products in `docs/`.
   - Update `README.md` to include new product support.

5. **Release**
   - Tag a release (e.g., `v2.0.0`) with new product support.
   - Publish updated documentation and examples.

---

## Timeline
| Phase | Duration | Status |
|-------|----------|--------|
| Phase 1: Architecture & Design | 4-6 weeks | Not Started |
| Phase 2: Existing Implementation Wrap-Up | 6-8 weeks | Not Started |
| Phase 3: New Product Implementation | 8-12 weeks | Not Started |

---

## Risks & Mitigations
- **Technical Debt**: Allocate time in Phase 1 to address audit findings.
- **Scope Creep**: Prioritize features based on product team feedback.
- **Testing Gaps**: Expand test coverage early to avoid regressions.

---

## Next Steps
1. Review and refine this roadmap with stakeholders.
2. Prioritize Phase 1 tasks and create GitHub/GitLab issues.
3. Begin Phase 1 implementation.