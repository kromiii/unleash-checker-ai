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

func (c *Client) ModifyCode(content string, flags []string) (string, error) {
	ctx := context.Background()
	messages := []ChatCompletionMessage{
		{
			Role: ChatMessageRoleSystem,
			Content: `
      あなたは優秀なプログラマーです。与えられたコードを修正する任務があります。
      以下の指示に従ってコードを修正してください：

      1. 指定された未使用のフラグをコードから削除してください。
      2. これらのフラグは長期間有効で安定して動作しているため、フラグの状態を参照する必要はありません。
      3. フラグを参照している部分については、フラグが有効な場合の処理を継続してください。
      4. 修正されたコードの差分のみを出力してください。
      5. 各差分は以下の形式で出力してください：
         行番号:修正後のコード
      6. 行を削除する場合は、その行番号に対応する出力を空にしてください。
      7. 新しい行を追加する場合は、_:（冒頭に追加）または+:（末尾に追加）を使用してください。
      8. マークダウンやコードブロックは使用しないでください。
      9. コメントは必要最小限にし、コードの変更に集中してください。
      10. 変更がない行は出力しないでください。
      11. 行番号は1から始まることに注意してください。

      修正後のコードは、元のコードと同じ機能を保持しつつ、より簡潔で効率的になるようにしてください。

      入力例:

      1:
      2:export function add(a: number, b: number): number {
      3: return 0;
      4:}
      5:

      出力例1: 単純な置換

      3: return a + b;

      出力例2: 複数行への書き換え

      3: // implement add
      3: return a + b;

      出力例3: 冒頭に追加

      _:import {x} from 'module';

      出力例4: 末尾に追加

      +:// This file is modified by Code Assistant

      NG:

      6: return a + b; // This line number does not existed
      `,
		},
		{
			Role:    ChatMessageRoleUser,
			Content: fmt.Sprintf("以下のコードを修正してください：\n%s\n\n削除するフラグ: %s", content, strings.Join(flags, ", ")),
		},
	}

	resp, err := c.CreateChatCompletion(ctx, ChatCompletionRequest{
		Model:    Model,
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
	Model                 = "gpt-4"
	ChatMessageRoleUser   = "user"
	ChatMessageRoleSystem = "system"
)
