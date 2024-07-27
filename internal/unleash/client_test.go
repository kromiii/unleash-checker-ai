package unleash

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetStaleFlags(t *testing.T) {
	// モックサーバーのセットアップ
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		response := FeatureFlagsResponse{
			Features: []FeatureFlag{
				{Name: "flag1", Type: "release", CreatedAt: time.Now().Add(-41 * 24 * time.Hour), Enabled: true, Stale: false},
				{Name: "flag2", Type: "experiment", CreatedAt: time.Now().Add(-30 * 24 * time.Hour), Enabled: true, Stale: false},
				{Name: "flag3", Type: "operational", CreatedAt: time.Now().Add(-8 * 24 * time.Hour), Enabled: true, Stale: true},
				{Name: "flag4", Type: "killSwitch", CreatedAt: time.Now().Add(-366 * 24 * time.Hour), Enabled: true, Stale: false},
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	// クライアントの作成
	client := NewClient(server.URL, "test-token", "test-project")

	// テストケース
	tests := []struct {
		name           string
		onlyStaleFlags bool
		expected       []string
	}{
		{"全ての古いフラグを取得", false, []string{"flag1", "flag3", "flag4"}},
		{"Staleフラグのみを取得", true, []string{"flag3"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flags, err := client.GetStaleFlags(tc.onlyStaleFlags)
			if err != nil {
				t.Fatalf("GetStaleFlags failed: %v", err)
			}

			if len(flags) != len(tc.expected) {
				t.Errorf("Expected %d flags, got %d", len(tc.expected), len(flags))
			}

			for i, flag := range flags {
				if flag != tc.expected[i] {
					t.Errorf("Expected flag %s, got %s", tc.expected[i], flag)
				}
			}
		})
	}
}

func TestGetExpectedLifetime(t *testing.T) {
	tests := []struct {
		flagType string
		expected time.Duration
	}{
		{"release", 40 * 24 * time.Hour},
		{"experiment", 40 * 24 * time.Hour},
		{"operational", 7 * 24 * time.Hour},
		{"killSwitch", 365 * 24 * time.Hour},
		{"permission", 365 * 24 * time.Hour},
		{"unknown", 30 * 24 * time.Hour},
	}

	for _, tc := range tests {
		t.Run(tc.flagType, func(t *testing.T) {
			result := getExpectedLifetime(tc.flagType)
			if result != tc.expected {
				t.Errorf("Expected %v for %s, got %v", tc.expected, tc.flagType, result)
			}
		})
	}
}