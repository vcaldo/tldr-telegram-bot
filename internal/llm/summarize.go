package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const ollamaAPIURL = "http://ollama:11434/api/generate"

// Summarize sends a request to the Ollama LLM server to summarize the provided text based on the specified language.
func Summarize(text string, lang string) (string, error) {
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		return "", fmt.Errorf("OLLAMA_MODEL environment variable is not set")
	}

	prompt := constructPrompt(text, lang)

	requestBody, err := json.Marshal(map[string]string{
		"model":  model,
		"prompt": prompt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(ollamaAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response from Ollama API: %s", resp.Status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response from Ollama API: %v", err)
	}

	summary, ok := response["summary"].(string)
	if !ok {
		return "", fmt.Errorf("summary not found in response")
	}

	return summary, nil
}

// constructPrompt creates a prompt for the LLM based on the specified language.
func constructPrompt(text string, lang string) string {
	switch lang {
	case "pt":
		return fmt.Sprintf("Resuma a seguinte Telegram chat in Portuguese:\n%s", text)
	case "en":
		return fmt.Sprintf("Summarize the following Telegram chat in English:\n%s", text)
	case "es":
		return fmt.Sprintf("Resume el siguiente Telegram chat en español:\n%s", text)
	default:
		return text // Fallback to plain text if language is unknown
	}
}