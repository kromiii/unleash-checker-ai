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
}

type FeatureFlag struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

type FeatureFlagsResponse struct {
	Features []FeatureFlag `json:"features"`
}

func NewUnleashClient(baseURL, apiToken string) *UnleashClient {
	return &UnleashClient{
		BaseURL:  baseURL,
		APIToken: apiToken,
	}
}

func (c *UnleashClient) GetStaleFlags() ([]string, error) {
	url := fmt.Sprintf("%s/api/admin/features", c.BaseURL)

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

	return getStaleFlags(featureFlagsResp.Features), nil
}

func getStaleFlags(flags []FeatureFlag) []string {
	var staleFlags []string
	now := time.Now()

	for _, flag := range flags {
		lifetime := getExpectedLifetime(flag.Type)
		if now.Sub(flag.CreatedAt) > lifetime {
			staleFlags = append(staleFlags, flag.Name)
		}
	}

	return staleFlags
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
