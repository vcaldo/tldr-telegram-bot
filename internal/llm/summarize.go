package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Summarize sends a request to the Ollama LLM server and waits until done is true.
func Summarize(text string, lang string) (string, error) {
	ollamaAPIURL := os.Getenv("OLLAMA_API_URL")
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

	var summary string
	scanner := bufio.NewScanner(resp.Body)
	// Opcional: definir um timeout se necessário
	done := false
	for scanner.Scan() {
		line := scanner.Bytes()
		var partial map[string]interface{}
		if err := json.Unmarshal(line, &partial); err != nil {
			log.Printf("failed to unmarshal chunk: %v", err)
			continue
		}
		log.Printf("Partial response from Ollama API: %v", partial)

		// Atualiza o summary se houver "response"
		if r, ok := partial["response"].(string); ok {
			summary = r
		}

		// Verifica se o processamento foi finalizado.
		if d, ok := partial["done"].(bool); ok && d {
			done = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Se não encontrou summary ou processing não terminou, podemos esperar um pouquinho
	// ou retornar erro, conforme sua estratégia.
	if !done {
		// Pode aguardar um tempo adicional aqui se preferir, por exemplo:
		time.Sleep(2 * time.Second)
		if summary == "" {
			return "", fmt.Errorf("summary not found in response after waiting")
		}
	}

	return summary, nil
}

// constructPrompt creates a prompt for the LLM based on the specified language.
func constructPrompt(text string, lang string) string {
	switch lang {
	case "pt":
		return fmt.Sprintf("Resuma a seguinte conversa do Telegram em Português:\n%s", text)
	case "en":
		return fmt.Sprintf("Summarize the following Telegram chat in English:\n%s", text)
	case "es":
		return fmt.Sprintf("Resume el siguiente Telegram chat en español:\n%s", text)
	default:
		return text // fallback para texto puro
	}
}
