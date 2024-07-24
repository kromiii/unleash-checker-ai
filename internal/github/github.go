package github

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	owner  string
	repo   string
}

func NewClient(token, owner, repo, baseURL string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// GHES のエンドポイントを設定
	if baseURL != "" {
		baseEndpoint, _ := url.Parse(baseURL + "/api/v3/")
		uploadEndpoint, _ := url.Parse(baseURL + "/api/uploads/")
		client.BaseURL = baseEndpoint
		client.UploadURL = uploadEndpoint
	}

	return &Client{
		client: client,
		owner:  owner,
		repo:   repo,
	}
}

func (c *Client) GetDefaultBranch(ctx context.Context) (string, error) {
	repo, _, err := c.client.Repositories.Get(ctx, c.owner, c.repo)
	if err != nil {
		return "", fmt.Errorf("error getting repository: %v", err)
	}

	return repo.GetDefaultBranch(), nil
}

func (c *Client) CreateBranch(ctx context.Context, branch, baseBranch string) error {
	ref, _, err := c.client.Git.GetRef(ctx, c.owner, c.repo, "refs/heads/"+baseBranch)
	if err != nil {
		return fmt.Errorf("error getting base branch ref: %v", err)
	}

	newRef := &github.Reference{
		Ref:    github.String("refs/heads/" + branch),
		Object: &github.GitObject{SHA: ref.Object.SHA},
	}

	_, _, err = c.client.Git.CreateRef(ctx, c.owner, c.repo, newRef)
	if err != nil {
		return fmt.Errorf("error creating new branch: %v", err)
	}

	return nil
}

func (c *Client) CommitChanges(ctx context.Context, branch, message string, files []string) error {
	// Get the current commit SHA
	ref, _, err := c.client.Git.GetRef(ctx, c.owner, c.repo, "refs/heads/"+branch)
	if err != nil {
		return fmt.Errorf("error getting ref: %v", err)
	}

	// Create a tree with the updated files
	var entries []*github.TreeEntry
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", file, err)
		}

		entries = append(entries, &github.TreeEntry{
			Path:    github.String(file),
			Mode:    github.String("100644"),
			Type:    github.String("blob"),
			Content: github.String(string(content)),
		})
	}

	tree, _, err := c.client.Git.CreateTree(ctx, c.owner, c.repo, *ref.Object.SHA, entries)
	if err != nil {
		return fmt.Errorf("error creating tree: %v", err)
	}

	// Create a new commit
	parent, _, err := c.client.Repositories.GetCommit(ctx, c.owner, c.repo, *ref.Object.SHA, nil)
	if err != nil {
		return fmt.Errorf("error getting parent commit: %v", err)
	}

	commit, _, err := c.client.Git.CreateCommit(ctx, c.owner, c.repo, &github.Commit{
		Message: github.String(message),
		Tree:    tree,
		Parents: []*github.Commit{{SHA: parent.SHA}},
	})
	if err != nil {
		return fmt.Errorf("error creating commit: %v", err)
	}

	// Update the reference
	_, _, err = c.client.Git.UpdateRef(ctx, c.owner, c.repo, &github.Reference{
		Ref:    github.String("refs/heads/" + branch),
		Object: &github.GitObject{SHA: commit.SHA},
	}, false)
	if err != nil {
		return fmt.Errorf("error updating ref: %v", err)
	}

	return nil
}

func (c *Client) CreatePullRequest(ctx context.Context, title, body, head, base string) (*github.PullRequest, error) {
	// Check if the branch exists
	_, _, err := c.client.Repositories.GetBranch(ctx, c.owner, c.repo, head, false)
	if err != nil {
		// Branch doesn't exist, create it
		ref, _, err := c.client.Git.GetRef(ctx, c.owner, c.repo, "refs/heads/"+base)
		if err != nil {
			return nil, fmt.Errorf("error getting base branch ref: %v", err)
		}

		newRef := &github.Reference{
			Ref:    github.String("refs/heads/" + head),
			Object: &github.GitObject{SHA: ref.Object.SHA},
		}

		_, _, err = c.client.Git.CreateRef(ctx, c.owner, c.repo, newRef)
		if err != nil {
			return nil, fmt.Errorf("error creating new branch: %v", err)
		}
	}

	// 既存のプルリクエストを検索
	opts := &github.PullRequestListOptions{
		State: "open",
		Head:  fmt.Sprintf("%s:%s", c.owner, head),
		Base:  base,
	}
	prs, _, err := c.client.PullRequests.List(ctx, c.owner, c.repo, opts)
	if err != nil {
		return nil, fmt.Errorf("error listing pull requests: %v", err)
	}

	// 該当するプルリクエストが見つかった場合、更新を行う
	if len(prs) > 0 {
		pr := prs[0] // 最初のプルリクエストを対象とする
		pr.Title = github.String(title)
		pr.Body = github.String(body)
		return pr, nil
	} else {
		// 該当するプルリクエストがない場合、新しいプルリクエストを作成
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
}
