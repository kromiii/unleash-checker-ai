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
    result := []string{}
    lineIndex := 0

    for _, diffLine := range diffLines {
        parts := strings.SplitN(diffLine, ": ", 2)
        if len(parts) != 2 {
            continue
        }

        if parts[0] == "_" {
            // 冒頭に新しい行を追加
            result = append(result, parts[1])
        } else if parts[0] == "+" {
            // 末尾に新しい行を追加
            result = append(result, parts[1])
        } else {
            lineNum, err := strconv.Atoi(parts[0])
            if err != nil {
                continue
            }

            // 変更されていない行を追加
            for lineIndex < lineNum-1 {
                result = append(result, lines[lineIndex])
                lineIndex++
            }

            if lineNum > 0 && lineNum <= len(lines) {
                if parts[1] == "" {
                    // 行を削除（何もしない）
                } else {
                    // 行を置換または修正（行番号を保持）
                    result = append(result, fmt.Sprintf("%d: %s", lineNum, parts[1]))
                }
                lineIndex = lineNum
            }
        }
    }

    // 残りの変更されていない行を追加
    for lineIndex < len(lines) {
        result = append(result, lines[lineIndex])
        lineIndex++
    }

    return strings.Join(result, "\n"), nil
}
