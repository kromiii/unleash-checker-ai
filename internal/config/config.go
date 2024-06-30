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
}

func Load() (*Config, error) {
    endpoint := os.Getenv("UNLEASH_API_ENDPOINT")
    token := os.Getenv("UNLEASH_API_TOKEN")
    projectID := os.Getenv("UNLEASH_PROJECT_ID")
    openaiKey := os.Getenv("OPENAI_API_KEY")

    if endpoint == "" || token == "" || openaiKey == "" || projectID == "" {
        return nil, errors.New("missing required environment variables")
    }

    return &Config{
        UnleashAPIEndpoint: endpoint,
        UnleashAPIToken:    token,
        ProjectID:          projectID,
        OpenAIAPIKey:       openaiKey,
    }, nil
}
