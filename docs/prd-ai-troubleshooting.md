# PRD: AI-driven Troubleshooting for Launchpad

## Overview
Launchpad performs complex cluster management tasks. When `launchpad apply` or other commands fail, users often need help diagnosing the root cause. This feature integrates AI to analyze failures and provide actionable troubleshooting advice.

## Goals
1. Automatically collect relevant context (logs, config, state) upon command failure.
2. Provide an AI-powered summary of the failure and potential fixes.
3. Integrate with the CLI via an `--ai-troubleshoot` flag or a reactive prompt.

## User Flow
1. User runs `launchpad apply --config cluster.yaml --ai-troubleshoot`.
2. A phase fails (e.g., MKE installation fails due to unreachable host).
3. Launchpad catches the error.
4. Launchpad collects:
   - The specific error message.
   - Relevant logs from `pkg/log`.
   - The `cluster.yaml` configuration (redacted of secrets).
   - Component state.
5. Launchpad sends this context to an AI provider.
6. Launchpad displays the AI's analysis and recommended next steps to the user.

## Technical Design
- **New Package**: `pkg/ai` to handle LLM communication.
- **Provider Interface**: Support for multiple providers (OpenAI, Azure, etc.) as mentioned in `mirantis-guidance.md`.
- **Context Collector**: A utility to gather redacted information for the AI prompt.
- **CLI Integration**: Modify `pkg/project/cli.go` and `cmd/bootstrap.go` to support AI flags.

## Success Criteria
- Users receive more specific advice than just "command failed".
- AI correctly identifies common issues like network timeouts, permission errors, or configuration mismatches.
- No sensitive data (API keys, passwords) is sent to the AI.

## Security & Compliance
- Follow `docs/ai/mirantis-guidance.md`.
- Ensure data redaction.
- Use environment variables for API keys.
