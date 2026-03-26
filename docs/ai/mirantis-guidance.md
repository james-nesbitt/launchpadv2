# Mirantis AI Guidance

*Corporate-wide policies and roles for AI tooling at Mirantis.*

---

## Roles
- **AI Compliance Officer**: Enforces Mirantis security/privacy policies for AI usage.
- **Vendor Liaison**: Manages relationships with LLM/API providers (e.g., OpenAI, Azure).
- **Security Reviewer**: Audits AI features for compliance.

---

## Maintenance
### Vendor Updates
- **Frequency**: Quarterly review of LLM/API integrations.
- **Process**:
  1. Test new models in staging.
  2. Update `pkg/ai/vendors/` with approved versions.
  3. Deprecate outdated models.

### Security
- **API Keys**: Encrypt and rotate quarterly.
- **Logging**: Audit all AI interactions for compliance.

---

## Compliance
- **Data Privacy**: No customer data in prompts or training.
- **Approval**: New AI features require security review.
- **Documentation**: Update `mirantis-guidance.md` for policy changes.