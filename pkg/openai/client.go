package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) ModifyCode(content string, unusedFlags []string) (string, error) {
	ctx := context.Background()
	messages := []ChatCompletionMessage{
		{
			Role:    ChatMessageRoleSystem,
			Content: "You are a code modification assistant. Remove the specified unused flags from the given code.",
		},
		{
			Role:    ChatMessageRoleUser,
			Content: fmt.Sprintf("Remove the following unused flags from this code:\n%s\n\nUnused flags: %s", content, strings.Join(unusedFlags, ", ")),
		},
	}

	resp, err := c.CreateChatCompletion(ctx, ChatCompletionRequest{
		Model:    GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", strings.NewReader(string(jsonBody)))
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ChatCompletionResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return ChatCompletionResponse{}, err
	}

	return response, nil
}

type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

const (
	GPT3Dot5Turbo         = "gpt-3.5-turbo"
	ChatMessageRoleUser   = "user"
	ChatMessageRoleSystem = "system"
)
