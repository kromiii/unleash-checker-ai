package unleash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UnleashClient struct {
	BaseURL  string
	APIToken string
	ProjectID string
}

type FeatureFlag struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	Enabled   bool      `json:"enabled"`
	Stale     bool      `json:"stale"`
}

type FeatureFlagsResponse struct {
	Features []FeatureFlag `json:"features"`
}

func NewUnleashClient(baseURL, apiToken string, projectID string) *UnleashClient {
	return &UnleashClient{
		BaseURL:  baseURL,
		APIToken: apiToken,
		ProjectID: projectID,
	}
}

func (c *UnleashClient) GetUnusedAndStaleFlags() ([]string, error) {
	url := fmt.Sprintf("%s/admin/projects/%s/features", c.BaseURL, c.ProjectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.APIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var featureFlagsResp FeatureFlagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&featureFlagsResp); err != nil {
		return nil, err
	}

	return getUnusedAndStaleFlags(featureFlagsResp.Features), nil
}

func getUnusedAndStaleFlags(flags []FeatureFlag) []string {
	var unusedAndStaleFlags []string
	now := time.Now()

	for _, flag := range flags {
		lifetime := getExpectedLifetime(flag.Type)
		isStale := flag.Stale || now.Sub(flag.CreatedAt) > lifetime

		if isStale {
			unusedAndStaleFlags = append(unusedAndStaleFlags, flag.Name)
		}
	}

	return unusedAndStaleFlags
}

func getExpectedLifetime(flagType string) time.Duration {
	switch flagType {
	case "release", "experiment":
		return 40 * 24 * time.Hour // 40 days
	case "operational":
		return 7 * 24 * time.Hour // 7 days
	case "killSwitch", "permission":
		return 365 * 24 * time.Hour // 1 year (as these are expected to be permanent)
	default:
		return 30 * 24 * time.Hour // Default to 30 days
	}
}
