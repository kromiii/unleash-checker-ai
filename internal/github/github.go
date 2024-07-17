package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	owner  string
	repo   string
}

func NewClient(token, owner, repo string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
		owner:  owner,
		repo:   repo,
	}
}

func (c *Client) CreatePullRequest(ctx context.Context, title, body, head, base string) (*github.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Body:  github.String(body),
		Head:  github.String(head),
		Base:  github.String(base),
	}

	pr, _, err := c.client.PullRequests.Create(ctx, c.owner, c.repo, newPR)
	if err != nil {
		return nil, fmt.Errorf("error creating pull request: %v", err)
	}

	return pr, nil
}
