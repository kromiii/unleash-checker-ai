package config

import (
	"errors"
	"os"
)

type Config struct {
	UnleashAPIEndpoint string
	UnleashAPIToken    string
	ProjectID          string
	OpenAIAPIKey       string
	GitHubToken        string
	GitHubOwner        string
	GitHubRepo         string
	GitHubBaseURL			string
}

func Load() (*Config, error) {
	endpoint := os.Getenv("UNLEASH_API_ENDPOINT")
	token := os.Getenv("UNLEASH_API_TOKEN")
	projectID := os.Getenv("UNLEASH_PROJECT_ID")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	githubToken := os.Getenv("GITHUB_TOKEN")
	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")
	githubBaseURL := os.Getenv("GITHUB_BASE_URL")

	if endpoint == "" || token == "" || openaiKey == "" || projectID == "" || githubToken == "" || githubOwner == "" || githubRepo == "" {
		return nil, errors.New("missing required environment variables")
	}

	return &Config{
		UnleashAPIEndpoint: endpoint,
		UnleashAPIToken:    token,
		ProjectID:          projectID,
		OpenAIAPIKey:       openaiKey,
		GitHubToken:        githubToken,
		GitHubOwner:        githubOwner,
		GitHubRepo:         githubRepo,
		GitHubBaseURL:      githubBaseURL,
	}, nil
}
