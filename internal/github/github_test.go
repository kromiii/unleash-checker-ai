package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-github/v38/github"
)

func TestNewClient(t *testing.T) {
	client := NewClient("token", "owner", "repo", "")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.owner != "owner" {
		t.Errorf("Expected owner to be 'owner', got %s", client.owner)
	}
	if client.repo != "repo" {
		t.Errorf("Expected repo to be 'repo', got %s", client.repo)
	}
}

func TestGetDefaultBranch(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/repos/owner/repo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"default_branch": "main"}`))
	})

	client := &Client{
		client: github.NewClient(nil),
		owner:  "owner",
		repo:   "repo",
	}
	baseURL, _ := url.Parse(server.URL + "/")
	client.client.BaseURL = baseURL
	branch, err := client.GetDefaultBranch(context.Background())
	if err != nil {
		t.Fatalf("GetDefaultBranch returned error: %v", err)
	}
	if branch != "main" {
		t.Errorf("Expected default branch to be 'main', got %s", branch)
	}
}

func TestCreateBranch(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/repos/owner/repo/git/ref/heads/main", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"object": {"sha": "abcdef1234567890"}}`))
	})

	mux.HandleFunc("/repos/owner/repo/git/refs", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	client := &Client{
		client: github.NewClient(nil),
		owner:  "owner",
		repo:   "repo",
	}
	baseURL, _ := url.Parse(server.URL + "/")
	client.client.BaseURL = baseURL

	err := client.CreateBranch(context.Background(), "new-branch", "main")
	if err != nil {
		t.Fatalf("CreateBranch returned error: %v", err)
	}
}

// 他のメソッドのテストも同様に実装してください
func TestCommitChanges(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/repos/owner/repo/git/refs/heads/test-branch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"object": {"sha": "abcdef1234567890"}}`))
		} else if r.Method == http.MethodPatch {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ref": "refs/heads/test-branch", "object": {"sha": "newcommitsha"}}`))
		} else {
			http.Error(w, "不正なメソッド", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/repos/owner/repo/git/trees", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"sha": "newtreesha"}`))
	})

	mux.HandleFunc("/repos/owner/repo/commits/abcdef1234567890", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sha": "parentcommitsha"}`))
	})

	mux.HandleFunc("/repos/owner/repo/git/commits", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"sha": "newcommitsha"}`))
	})

	client := &Client{
		client: github.NewClient(nil),
		owner:  "owner",
		repo:   "repo",
	}
	baseURL, _ := url.Parse(server.URL + "/")
	client.client.BaseURL = baseURL

	// テスト用の一時ファイルを作成
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("一時ファイルの作成に失敗しました: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("テストコンテンツ")
	if err != nil {
		t.Fatalf("一時ファイルへの書き込みに失敗しました: %v", err)
	}
	tempFile.Close()

	err = client.CommitChanges(context.Background(), "test-branch", "テストコミット", []string{tempFile.Name()})
	if err != nil {
		t.Fatalf("CommitChangesがエラーを返しました: %v", err)
	}
}

func TestCreatePullRequest(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/repos/owner/repo/branches/new-branch", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "new-branch"}`))
	})

	mux.HandleFunc("/repos/owner/repo/pulls", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"number": 1, "title": "テストPR", "body": "PRの本文"}`))
		}
	})

	client := &Client{
		client: github.NewClient(nil),
		owner:  "owner",
		repo:   "repo",
	}
	baseURL, _ := url.Parse(server.URL + "/")
	client.client.BaseURL = baseURL

	pr, err := client.CreatePullRequest(context.Background(), "テストPR", "PRの本文", "new-branch", "main")
	if err != nil {
		t.Fatalf("CreatePullRequestがエラーを返しました: %v", err)
	}

	if pr.GetNumber() != 1 {
		t.Errorf("期待されるPR番号は1ですが、%dが返されました", pr.GetNumber())
	}
	if pr.GetTitle() != "テストPR" {
		t.Errorf("期待されるPRタイトルは'テストPR'ですが、'%s'が返されました", pr.GetTitle())
	}
	if pr.GetBody() != "PRの本文" {
		t.Errorf("期待されるPR本文は'PRの本文'ですが、'%s'が返されました", pr.GetBody())
	}
}