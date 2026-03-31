package project

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/ai"
)

// Troubleshoot command failure with AI.
func (p *Project) Troubleshoot(ctx context.Context, command string, runErr error) error {
	slog.InfoContext(ctx, "Command failed. Attempting AI troubleshooting...")

	// Gather context
	req := ai.TroubleshootingRequest{
		Command: command,
		Error:   runErr.Error(),
		// We can gather more context here:
		// Config:  p.RedactedConfig(),
		// Logs:    log.RecentLogs(100),
		// State:   p.DebugState(),
	}

	// For now, use the mock provider
	// In the future, this could be configurable from the project config
	provider := ai.NewDefaultProvider()

	resp, tserr := provider.Troubleshoot(ctx, req)
	if tserr != nil {
		return fmt.Errorf("AI troubleshooting failed: %w", tserr)
	}

	fmt.Printf("\n--- AI TROUBLESHOOTING ANALYSIS ---\n\n")
	fmt.Printf("Summary: %s\n\n", resp.Summary)
	fmt.Printf("Analysis: %s\n\n", resp.Analysis)
	fmt.Printf("Recommendations:\n")
	for _, r := range resp.Recommendations {
		fmt.Printf(" - %s\n", r)
	}
	fmt.Printf("\n------------------------------------\n")

	return nil
}
