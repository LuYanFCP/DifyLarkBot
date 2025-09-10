package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

type ChatRequest struct {
	Inputs         map[string]any `json:"inputs"`
	Query          string                 `json:"query"`
	ResponseMode   string                 `json:"response_mode"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	User           string                 `json:"user"`
}

type ChatResponse struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Answer         string `json:"answer"`
}

func NewClient(apiKey, baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://api.dify.ai"
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *Client) Chat(message, userID string) (*ChatResponse, error) {
	url := c.baseURL + "/v1/chat-messages"
	
	requestBody := ChatRequest{
		Inputs:         make(map[string]any),
		Query:          message,
		ResponseMode:   "blocking",
		User:           userID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("dify API status %d: %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}