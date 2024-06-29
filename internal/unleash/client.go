package unleash

import (
	"github.com/Unleash/unleash-client-go/v4"
)

type Client struct {
	client *unleash.Client
}

func NewClient(url, apiToken string) (*Client, error) {
	client, err := unleash.NewClient(
		unleash.WithUrl(url),
		unleash.WithCustomHeaders(map[string][]string{"Authorization": {apiToken}}), // 修正箇所
		unleash.WithAppName("unleash-checker-ai"),
	)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) GetUnusedFlags() ([]string, error) {
	// この部分は Unleash API の実際の仕様に合わせて実装する必要があります
	// ここでは、ダミーの実装を返しています
	return []string{"flag_1", "flag_2"}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
