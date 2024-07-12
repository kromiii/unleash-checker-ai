package modifier

import (
	"os"
	"regexp"
	"strings"

	"github.com/kromiii/unleash-checker-ai/pkg/openai"
)

type Modifier struct {
	openaiClient *openai.Client
}

func NewModifier(apiKey string) *Modifier {
	return &Modifier{
		openaiClient: openai.NewClient(apiKey),
	}
}

func (m *Modifier) ModifyFile(filePath string, unusedFlags []string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	modifiedContent := string(content)
	lines := strings.Split(modifiedContent, "\n")
	removedFlags := []string{}
	for i, line := range lines {
		if found, matchedFlag := findMatchedFlag(line, unusedFlags); found {
			lines[i] = "// This feature flag is stale and can be removed: " + matchedFlag + "\n" + line
			removedFlags = append(removedFlags, matchedFlag)
		}
	}

	modifiedContent, err = m.openaiClient.ModifyCode(strings.Join(lines, "\n"), "")
	if err != nil {
		return nil, err
	}
	
	if !strings.HasSuffix(modifiedContent, "\n") {
		modifiedContent += "\n"
	}

	err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
	return removedFlags, err
}

func findMatchedFlag(s string, flags []string) (bool, string) {
	pattern := strings.Join(flags, "|")
	re, err := regexp.Compile(pattern)
	if err != nil {
			return false, ""
	}
	matched := re.FindString(s)
	return matched != "", matched
}
