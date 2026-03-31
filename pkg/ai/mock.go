package ai

import (
	"context"
	"fmt"
)

// MockProvider is a mock AI troubleshooting provider for testing.
type MockProvider struct{}

// Troubleshoot analyze the request and return a mock response.
func (m *MockProvider) Troubleshoot(ctx context.Context, req TroubleshootingRequest) (TroubleshootingResponse, error) {
	return TroubleshootingResponse{
		Summary:  fmt.Sprintf("Mock AI: Analyzed failure in command '%s'", req.Command),
		Analysis: "Your command failed because of a simulated error. This mock AI provider suggests you check your environment and try again.",
		Recommendations: []string{
			"Verify that all cluster nodes are reachable.",
			"Check for any networking errors in the logs.",
			"Ensure you have sufficient permissions to perform the operation.",
		},
		Confidence: 0.95,
	}, nil
}
