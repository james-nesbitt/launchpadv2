# Code Guidelines for AI Agents and Developers

This document outlines the **coding standards** and **best practices** for this repository. AI agents and developers should follow these guidelines to ensure consistency, readability, and maintainability.

---

## 1. General Principles
### 1.1 Consistency
- Follow existing patterns in the codebase.
- Use the same naming conventions, formatting, and structure as the rest of the project.

### 1.2 Readability
- Write code for humans first, machines second.
- Use descriptive variable and function names.
- Add comments for complex logic or non-obvious decisions.

### 1.3 Maintainability
- Keep functions small and focused (single responsibility principle).
- Avoid hardcoding values; use constants or configuration.
- Ensure backward compatibility where possible.

---

## 2. Go-Specific Guidelines
### 2.1 Formatting
- Use `gofmt` to format code automatically.
  ```bash
  gofmt -w .
  ```
- Use `golangci-lint` to catch style issues.
  ```bash
  make lint
  ```

### 2.2 Naming Conventions
| Type | Convention | Example |
|------|------------|---------|
| Packages | Short, lowercase, no underscores | `engine`, `config` |
| Variables | camelCase | `clusterName`, `isValid` |
| Constants | UPPER_SNAKE_CASE | `DEFAULT_TIMEOUT` |
| Functions | camelCase | `ApplyCluster()`, `ValidateConfig()` |
| Interfaces | Single-method: `-er` suffix<br>Multi-method: descriptive name | `Driver`, `ConfigValidator` |
| Files | lowercase_with_underscores.go | `driver_interface.go` |

### 2.3 Error Handling
- Use `errors.Wrap()` or `fmt.Errorf()` to add context to errors.
- Return errors explicitly; avoid swallowing them.
- Example:
  ```go
  if err := someFunction(); err != nil {
      return fmt.Errorf("failed to do X: %w", err)
  }
  ```

### 2.4 Testing
- Write unit tests for all exported functions.
- Use table-driven tests for complex logic.
- Example:
  ```go
  func TestValidateConfig(t *testing.T) {
      tests := []struct {
          name     string
          config   Config
          wantErr  bool
      }{
          {"valid config", Config{...}, false},
          {"invalid config", Config{...}, true},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              if err := ValidateConfig(tt.config); (err != nil) != tt.wantErr {
                  t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
              }
          })
      }
  }
  ```

### 2.5 Logging
- Use the standard `log` package for logging.
- Include context in logs (e.g., `log.Printf("Applying cluster %s: %v", clusterName, err)`).
- Use `log.Fatal()` only for unrecoverable errors.

---

## 3. Documentation
### 3.1 Godoc Comments
- Add Godoc comments for all exported functions, types, and packages.
- Example:
  ```go
  // ApplyCluster applies the configuration to create or update a cluster.
  // It returns an error if the configuration is invalid or the operation fails.
  func ApplyCluster(cfg Config) error {
      // ...
  }
  ```

### 3.2 Inline Comments
- Add comments for complex logic or non-obvious decisions.
- Avoid stating the obvious (e.g., `// Increment i`).
- Example:
  ```go
  // Retry the operation up to 3 times with exponential backoff.
  for i := 0; i < 3; i++ {
      if err := doOperation(); err == nil {
          break
      }
      time.Sleep(time.Duration(1<<i) * time.Second)
  }
  ```

### 3.3 AI Agent-Specific Documentation
- Update the `docs/ai/` folder when:
  - Adding new features or workflows.
  - Changing core interfaces or abstractions.
  - Introducing breaking changes.
- Ensure the `README.md` in `docs/ai/` is always up-to-date.

---

## 4. AI Agent Guidance
### 4.1 Contributing Code
- Follow the workflows in [`workflows.md`](./workflows.md).
- Ensure your changes align with these guidelines.
- Run `make test` and `make lint` before submitting a PR.

### 4.2 Reviewing Code
- Check for consistency with existing patterns.
- Verify that tests and documentation are updated.
- Ensure error handling is robust.

### 4.3 Suggesting Improvements
- Propose changes that align with the project's goals (e.g., maintainability, readability).
- Avoid suggestions that introduce complexity without clear benefits.
- Reference these guidelines when providing feedback.

---

## Next Steps
1. Review the [Workflows](./workflows.md) for task-specific guidance.
2. Check the [README](./README.md) for additional context.