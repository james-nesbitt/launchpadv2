package ai

import (
	"os"
)

// NewDefaultProvider return a new AI troubleshooting provider based on environment.
func NewDefaultProvider() Provider {
	// If OPENAI_API_KEY is present, use OpenAI
	// In the future, this can be more sophisticated
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		return NewOpenAIProvider(apiKey)
	}

	// Default to mock provider for now
	return &MockProvider{}
}
