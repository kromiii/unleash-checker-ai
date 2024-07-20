package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kromiii/unleash-checker-ai/internal/config"
	"github.com/kromiii/unleash-checker-ai/internal/finder"
	"github.com/kromiii/unleash-checker-ai/internal/github"
	"github.com/kromiii/unleash-checker-ai/internal/report"
	"github.com/kromiii/unleash-checker-ai/internal/unleash"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: unleash-checker-ai <folder>")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	onlyStaleFlag := flag.Bool("only-stale", false, "Ignore potentially stale flags")
	flag.Parse()

	client := unleash.NewClient(cfg.UnleashAPIEndpoint, cfg.UnleashAPIToken, cfg.ProjectID)
	onlyStaleFlags := *onlyStaleFlag
	staleFlags, err := client.GetStaleFlags(onlyStaleFlags)
	if err != nil {
		fmt.Printf("Error getting stale flags: %v\n", err)
		return
	}

	targetFolder := os.Args[1]
	changedFiles, removedFlags, err := finder.FindAndReplaceFlags(targetFolder, staleFlags, cfg.OpenAIAPIKey)
	if err != nil {
		fmt.Printf("Error finding affected files: %v\n", err)
		os.Exit(1)
	}

	summary := report.CreateSummary(staleFlags, removedFlags)

	// chanedFilesに変更がない場合は終了
	if len(changedFiles) == 0 {
		fmt.Println("No changes required")
		return
	}

	// Create GitHub client
	githubClient := github.NewClient(cfg.GitHubToken, cfg.GitHubOwner, cfg.GitHubRepo)

	// Commit changes
	ctx := context.Background()
	branchName := "unleash-checker-updates"
	err = githubClient.CommitChanges(ctx, branchName, "Update stale Unleash flags", changedFiles)
	if err != nil {
		fmt.Printf("Error committing changes: %v\n", err)
		os.Exit(1)
	}

	// Create pull request
	pr, err := githubClient.CreatePullRequest(ctx, "Update stale Unleash flags", summary, branchName, "main")
	if err != nil {
		fmt.Printf("Error creating pull request: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Pull request created: %s\n", pr.GetHTMLURL())
}
