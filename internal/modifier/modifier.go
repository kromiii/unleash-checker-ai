package modifier

import (
    "os"
    "strconv"
    "strings"
    "fmt"

    "github.com/kromiii/unleash-checker-ai/pkg/openai"
)

type OpenAIClientInterface interface {
    ModifyCode(content string, flags []string) (string, error)
}

type Modifier struct {
    openaiClient OpenAIClientInterface
}

func NewModifier(apiKey string) *Modifier {
    return &Modifier{
        openaiClient: openai.NewClient(apiKey),
    }
}

func (m *Modifier) ModifyFile(filePath string, unusedFlags []string) error {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    numberedContent := addLineNumbers(string(content))
    diff, err := m.openaiClient.ModifyCode(numberedContent, unusedFlags)
	fmt.Println(diff)
    if err != nil {
        return err
    }

    modifiedContent, err := applyDiff(numberedContent, diff)
    if err != nil {
        return err
    }

    finalContent := stripLineNumbers(modifiedContent)
    err = os.WriteFile(filePath, []byte(finalContent), 0644)
    return err
}

func addLineNumbers(content string) string {
    lines := strings.Split(content, "\n")
    numberedLines := make([]string, len(lines))
    for i, line := range lines {
        numberedLines[i] = fmt.Sprintf("%d: %s", i+1, line)
    }
    return strings.Join(numberedLines, "\n")
}

func stripLineNumbers(content string) string {
    lines := strings.Split(content, "\n")
    strippedLines := make([]string, len(lines))
    for i, line := range lines {
        parts := strings.SplitN(line, ": ", 2)
        if len(parts) == 2 {
            strippedLines[i] = parts[1]
        } else {
            strippedLines[i] = line
        }
    }
    return strings.Join(strippedLines, "\n")
}

func applyDiff(originalContent, diff string) (string, error) {
    lines := strings.Split(originalContent, "\n")
    diffLines := strings.Split(diff, "\n")

    for _, diffLine := range diffLines {
        parts := strings.SplitN(diffLine, ": ", 2)
        if len(parts) != 2 {
            continue
        }

        lineNum, err := strconv.Atoi(parts[0])
        if err != nil {
            if parts[0] == "*" {
                lines = append([]string{parts[1]}, lines...)
            } else if parts[0] == "+" {
                lines = append(lines, parts[1])
            }
            continue
        }

        if lineNum > 0 && lineNum <= len(lines) {
            if parts[1] == "" {
                lines = append(lines[:lineNum-1], lines[lineNum:]...)
            } else {
                lines[lineNum-1] = fmt.Sprintf("%d: %s", lineNum, parts[1])
            }
        }
    }

    return strings.Join(lines, "\n"), nil
}
