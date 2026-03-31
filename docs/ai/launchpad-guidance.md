# Launchpad AI Guidance (Future)

This document outlines **planned AI integrations** for Launchpad. No AI features are currently implemented.

---

## Planned Features
### 1. AI Troubleshooting
- **Command**: `launchpad discover --ai-debug` (not yet implemented).
- **Purpose**: Analyze cluster issues using AI.

### 2. AI-Assisted Effort Initialization
- **Purpose**: Guide users through starting a new development effort.
- **Steps**: When initiating a new effort, the AI will:
  1. **Assign a short label** for the effort (e.g., `feature/ai-debug`, `update/dependencies`).
  2. **Create a new branch** for the effort.
  3. **Collect details** for a new PRD (Product Requirements Document).
  4. **Review documentation**, including AI guidance and design docs, to ensure alignment with project goals.

---

## Roles (Future)
- **AI Feature Developer**: Implement AI-assisted CLI commands.
- **Code Reviewer**: Validate AI prompts for security and accuracy.

---

## Maintenance (Future)
- Update `pkg/ai/prompts/` when new models or components are adopted.