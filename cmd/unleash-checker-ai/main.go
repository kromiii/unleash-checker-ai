package main

import (
	"fmt"
	"os"
	"flag"

	"github.com/kromiii/unleash-checker-ai/internal/config"
	"github.com/kromiii/unleash-checker-ai/internal/finder"
	"github.com/kromiii/unleash-checker-ai/internal/modifier"
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
	unusedFlags, err := client.GetStaleFlags(onlyStaleFlags)
	if err != nil {
		fmt.Printf("Error getting stale flags: %v\n", err)
		return
	}

	fmt.Println("Stale or potentially stale flags:")
	for _, flag := range unusedFlags {
		fmt.Printf(" - %s\n", flag)
	}

	targetFolder := os.Args[1]
	affectedFiles, err := finder.FindAffectedFiles(targetFolder, unusedFlags)
	if err != nil {
		fmt.Printf("Error finding affected files: %v\n", err)
		os.Exit(1)
	}

	if len(affectedFiles) == 0 {
		fmt.Println("No files are affected")
		return
	}

	fmt.Println("These flags are used in:")
	for _, file := range affectedFiles {
		fmt.Printf(" - %s\n", file)
	}

	fmt.Println("Removing unused flags by LLM...")

	modifier := modifier.NewModifier(cfg.OpenAIAPIKey)
	for _, file := range affectedFiles {
		err := modifier.ModifyFile(file, unusedFlags)
		if err != nil {
			fmt.Printf("Error modifying file %s: %v\n", file, err)
			continue
		}
	}

	fmt.Println("Done!")
}
