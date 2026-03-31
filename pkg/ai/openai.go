package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	OpenAIModel    = "gpt-4o" // Suggested default
	OpenAIEndpoint = "https://api.openai.com/v1/chat/completions"
)

type OpenAIProvider struct {
	apiKey string
	model  string
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = OpenAIModel
	}
	return &OpenAIProvider{
		apiKey: apiKey,
		model:  model,
	}
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
}

func (p *OpenAIProvider) Troubleshoot(ctx context.Context, req TroubleshootingRequest) (TroubleshootingResponse, error) {
	prompt := fmt.Sprintf(`Analyze the following command failure and provide troubleshooting advice.
Command: %s
Error: %s
Config: %s
Logs: %s

Respond in JSON format with fields "summary", "analysis", "recommendations" (string array), and "confidence" (0-1).`,
		req.Command, req.Error, req.Config, req.Logs)

	apiReq := openAIRequest{
		Model: p.model,
		Messages: []openAIMessage{
			{Role: "system", Content: "You are an expert systems engineer and SRE specializing in Mirantis products and Kubernetes."},
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return TroubleshootingResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", OpenAIEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return TroubleshootingResponse{}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return TroubleshootingResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TroubleshootingResponse{}, fmt.Errorf("OpenAI API error: %s", resp.Status)
	}

	var apiResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return TroubleshootingResponse{}, err
	}

	if len(apiResp.Choices) == 0 {
		return TroubleshootingResponse{}, fmt.Errorf("no response from OpenAI")
	}

	var tsResp TroubleshootingResponse
	if err := json.Unmarshal([]byte(apiResp.Choices[0].Message.Content), &tsResp); err != nil {
		// Fallback if not JSON
		return TroubleshootingResponse{
			Summary:  "AI Analysis",
			Analysis: apiResp.Choices[0].Message.Content,
		}, nil
	}

	return tsResp, nil
}
