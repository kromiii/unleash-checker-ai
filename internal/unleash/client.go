package unleash

import (
	"encoding/json"
	"fmt"
	"math"
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

func NewClient(baseURL, apiToken string, projectID string) *UnleashClient {
	return &UnleashClient{
		BaseURL:  baseURL,
		APIToken: apiToken,
		ProjectID: projectID,
	}
}

func (c *UnleashClient) GetStaleFlags(onlyStaleFlags bool) ([]string, error) {
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

	return getStaleFlags(featureFlagsResp.Features, onlyStaleFlags), nil
}

func getStaleFlags(flags []FeatureFlag, onlyStaleFlags bool) []string {
	var staleFlags []string
	now := time.Now()

	for _, flag := range flags {
		lifetime := getExpectedLifetime(flag.Type)
		if onlyStaleFlags {
			if flag.Stale{
				staleFlags = append(staleFlags, flag.Name)
			}
		} else {
			isStale := flag.Stale || now.Sub(flag.CreatedAt) > lifetime
			if isStale {
				staleFlags = append(staleFlags, flag.Name)
			}
		}
	}

	return staleFlags
}

func getExpectedLifetime(flagType string) time.Duration {
	switch flagType {
	case "release", "experiment":
		return 40 * 24 * time.Hour // 40日
	case "operational":
		return 7 * 24 * time.Hour // 7日
	case "killSwitch", "permission":
		return time.Duration(math.MaxInt64) // 実質的に永続
	default:
		return 30 * 24 * time.Hour // デフォルトは30日
	}
}
