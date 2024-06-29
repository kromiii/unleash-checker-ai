package main

import (
    "fmt"
    "os"

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

    unleashClient, err := unleash.NewClient(cfg.UnleashAPIEndpoint, cfg.UnleashAPIToken)
    if err != nil {
        fmt.Printf("Error initializing Unleash client: %v\n", err)
        os.Exit(1)
    }
    defer unleashClient.Close()

    unusedFlags, err := unleashClient.GetUnusedFlags()
    if err != nil {
        fmt.Printf("Error getting unused flags: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Unused flags:")
    for _, flag := range unusedFlags {
        fmt.Printf(" - %s\n", flag)
    }

    targetFolder := os.Args[1]
    affectedFiles, err := finder.FindAffectedFiles(targetFolder, unusedFlags)
    if err != nil {
        fmt.Printf("Error finding affected files: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("These flags are used in:")
    for _, file := range affectedFiles {
        fmt.Printf(" - %s\n", file)
    }

    modifier := modifier.NewModifier(cfg.OpenAIAPIKey)
    for _, file := range affectedFiles {
        err := modifier.ModifyFile(file, unusedFlags)
        if err != nil {
            fmt.Printf("Error modifying file %s: %v\n", file, err)
            continue
        }
    }

    fmt.Println("Unused flags are removed")
}
