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
			Content: `
			You are a code modification assistant. Remove the specified unused flags from the given code.
			These flags remain enabled for a long time and operate stably, so we don't need to refer to the flag status. 
			You should only output the modified code because we will overwrite the original file with the modified code.
			Please provide the response in plain text without using any markdown or code blocks.
			`,
		},
		{
			Role:    ChatMessageRoleUser,
			Content: fmt.Sprintf("Modify the following code:\n%s\n\nStale flags are: %s", content, strings.Join(unusedFlags, ", ")),
		},
	}

	resp, err := c.CreateChatCompletion(ctx, ChatCompletionRequest{
		Model:    GPT4Turbo,
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
	GPT4Turbo         = "gpt-4-turbo"
	ChatMessageRoleUser   = "user"
	ChatMessageRoleSystem = "system"
)
