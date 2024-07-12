package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kromiii/unleash-checker-ai/internal/config"
	"github.com/kromiii/unleash-checker-ai/internal/finder"
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
	removedFlags, err := finder.FindAndReplaceFlags(targetFolder, staleFlags, cfg.OpenAIAPIKey)
	if err != nil {
		fmt.Printf("Error finding affected files: %v\n", err)
		os.Exit(1)
	}

	summary := report.CreateSummary(staleFlags, removedFlags)
	fmt.Println(summary)
}
