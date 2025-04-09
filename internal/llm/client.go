package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type OllamaClient struct {
	BaseURL    string
	HTTPClient *http.Client
	ModelName  string
}

func NewOllamaClient() *OllamaClient {
	return &OllamaClient{
		BaseURL:    os.Getenv("OLLAMA_API_URL"),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ModelName:  os.Getenv("OLLAMA_MODEL"),
	}
}

func (c *OllamaClient) GenerateSummary(prompt string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"model":  c.ModelName,
		"prompt": prompt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response from Ollama: %s", resp.Status)
	}

	var response struct {
		Summary string `json:"summary"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response from Ollama: %w", err)
	}

	return response.Summary, nil
}
