package docqa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is a wrapper around an LLM with schema capabilities.
type Client interface {
	// GetLLMResponse prompts the LLM and gets its text response, hopefully (not certainly) with the given schema.
	GetLLMResponse(systemPrompt string, userPrompt string, schema map[string]any) (string, error)
}

type openAIClient struct {
	key   string
	model string
}

// NewOpenAIClient creates a new client that communicates with the OpenAI API.
func NewOpenAIClient(key, model string) Client {
	return &openAIClient{
		key:   key,
		model: model,
	}
}

// GetLLMResponse implements [Client].
func (c *openAIClient) GetLLMResponse(systemPrompt string, userPrompt string, schema map[string]any) (string, error) {
	bodyMap := map[string]any{
		"model":           c.model,
		"temperature":     0.1,
		"response_format": wrapOpenAISchema(schema),
		"messages": []map[string]any{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
	}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respTyped := struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}{}
	err = json.Unmarshal(respBody, &respTyped)
	if err != nil || len(respTyped.Choices) == 0 || respTyped.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("failed to parse response: %s", string(respBody))
	}
	return respTyped.Choices[0].Message.Content, nil
}
