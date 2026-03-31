package ai

import "context"

// TroubleshootingRequest represents the context sent to the AI for analysis.
type TroubleshootingRequest struct {
	Command string `json:"command"`
	Error   string `json:"error"`
	Config  string `json:"config"` // Redacted cluster configuration
	Logs    string `json:"logs"`   // Relevant truncated logs
	State   any    `json:"state"`  // Current state of components
}

// TroubleshootingResponse represents the analysis and recommendations from the AI.
type TroubleshootingResponse struct {
	Summary         string   `json:"summary"`
	Analysis        string   `json:"analysis"`
	Recommendations []string `json:"recommendations"`
	Confidence      float64  `json:"confidence"`
}

// Provider defines the interface for an AI Troubleshooting service.
type Provider interface {
	Troubleshoot(ctx context.Context, req TroubleshootingRequest) (TroubleshootingResponse, error)
}
