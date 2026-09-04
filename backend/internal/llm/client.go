package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	apiURL string
	apiKey string
	model  string
	http   *http.Client
}

func NewClient(apiURL, apiKey, model string) *Client {
	return &Client{
		apiURL: apiURL,
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c.apiURL != ""
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *Client) Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody, err := json.Marshal(chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.4,
	})
	if err != nil {
		return "", fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	} else {
		req.Header.Set("Authorization", "")
		req.Header.Set("HTTP-Referer", "https://hermes-agent.nousresearch.com")
		req.Header.Set("X-Title", "Hearthside")
	}

	res, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("call llm: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errBody struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody.Error.Message != "" {
			return "", fmt.Errorf("llm returned status %d: %s", res.StatusCode, errBody.Error.Message)
		}
		return "", fmt.Errorf("llm returned status %d", res.StatusCode)
	}

	var parsed chatResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("llm returned no choices")
	}

	return parsed.Choices[0].Message.Content, nil
}

// ExtractJSON pulls a JSON array out of a model reply, tolerating markdown
// fences and surrounding prose. Falls back to wrapping a single JSON object
// in an array when the model ignores the array instruction.
func ExtractJSON(raw string) ([]byte, error) {
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start != -1 && end != -1 && end > start {
		return []byte(raw[start : end+1]), nil
	}

	objStart := strings.Index(raw, "{")
	objEnd := strings.LastIndex(raw, "}")
	if objStart != -1 && objEnd != -1 && objEnd > objStart {
		return []byte("[" + raw[objStart:objEnd+1] + "]"), nil
	}

	return nil, fmt.Errorf("no json found in model output")
}
