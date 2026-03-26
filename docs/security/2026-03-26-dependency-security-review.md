# Launchpad Security Review: Dependency Updates & Code Audit

**Date**: 2026-03-26
**Effort Branch**: `effort/update-vendor-libs-security`
**Reviewer**: Pi AI Assistant

---

## 1. Summary
This document captures the security findings from:
1. **Dependency updates** to latest stable versions.
2. **Code audit** for hardcoded secrets and sensitive data.

### Key Outcomes
| Area                  | Status       | Notes                                  |
|-----------------------|--------------|----------------------------------------|
| Dependency Updates    | ✅ Complete  | No vulnerabilities found via `govulncheck`. |
| Hardcoded Secrets     | ✅ None      | No production secrets in code.         |
| Input Validation      | ⚠️ Pending   | Requires manual review.                |
| Network Security      | ⚠️ Pending   | TLS/timeouts to be verified.           |

---

## 2. Dependency Updates
### Actions Taken
- Updated all Go dependencies to latest stable versions (see `go.mod`).
- Verified builds and tests pass (`go build ./...`, `go test ./...`).
- Scanned for vulnerabilities using `govulncheck`.

### Findings
- **No vulnerabilities** reported by `govulncheck`.
- **Breaking Changes**: None detected during build/test.

---

## 3. Code Audit for Hardcoded Secrets
### Methodology
- **Tool**: `grep -r "password\|api_key\|secret\|token" --include="*.go" .`
- **Scope**: All `.go` files in the repository.
- **Exclusions**: Test files (unit test values ignored).

### Findings
#### 3.1 Docker Swarm Tokens
- **Files**: `implementation/docker/docker.go`, `implementation/docker/dockerexec.go`
- **Context**: Tokens are **dynamically generated** for Swarm join operations.
- **Risk**: None. Tokens are ephemeral and not hardcoded.
- **Example**:
  ```go
  tcmd := []string{"swarm", "join-token", "-q"}  // Safe: Runtime-generated.
  ```

#### 3.2 K0s Join Tokens
- **Files**: `mirantis/product/k0s/hostplugin.go`, `mirantis/product/k0s/k0s.go`
- **Context**: Tokens are **dynamically generated** and protected by a mutex.
- **Risk**: None. Tokens are managed securely via `k0s token create`.
- **Example**:
  ```go
  "token",  // Safe: Runtime-generated via `k0s token create`.
  ```

#### 3.3 MKE3 Passwords
- **Files**: `mirantis/product/mke3/config.go`
- **Context**: `RegistryPassword` and `AdminPassword` are **user-provided** via YAML/CLI.
- **Risk**: None. No hardcoded values in production code.

#### 3.4 Test Files
- **Files**: `pkg/util/flags/flags_test.go`
- **Context**: Hardcoded test values (e.g., `--admin-password=barbar`).
- **Risk**: Low. Test-only data ignored per guidelines.

### Conclusion
- **No hardcoded secrets** found in production code.
- **Test files** contain placeholder values (ignored per scope).

---

## 4. Pending Security Reviews
### 4.1 Input Validation
- **Focus Areas**:
  - YAML parsing (`mirantis/config/v2_0/config.go`).
  - CLI flag validation (`cmd/launchpad/main.go`).
- **Next Steps**: Manual review for:
  - Sanitization of user inputs (e.g., hostnames, versions).
  - Validation of YAML/CLI schemas.

### 4.2 Network Security
- **Focus Areas**:
  - TLS usage in HTTP clients (e.g., `pkg/util/download/download.go`).
  - Default timeouts for network operations.
- **Next Steps**: Verify:
  - HTTPS/TLS for external connections (e.g., Helm, Docker).
  - Secure defaults for timeouts.

---

## 5. Recommendations
1. **Input Validation**:
   - Add schema validation for YAML/CLI inputs (e.g., regex for hostnames).
2. **Network Security**:
   - Enforce TLS 1.2+ for all external connections.
   - Set conservative timeouts (e.g., 30s for HTTP requests).
3. **Documentation**:
   - Update `docs/security.md` with findings and best practices.

---

## 6. Sign-off
- **Next Review**: Post-implementation of recommendations.
- **Owner**: Launchpad maintainers.