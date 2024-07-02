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

func (m *Modifier) ModifyFile(filePath string, unusedFlags []string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	modifiedContent := string(content)
	lines := strings.Split(modifiedContent, "\n")
	for i, line := range lines {
		if containsAny(line, unusedFlags) {
			modifiedLine, err := m.openaiClient.ModifyCode(line, unusedFlags)
			if err != nil {
				return err
			}

			lines[i] = modifiedLine
		}
	}

	modifiedContent = strings.Join(lines, "\n")

	return os.WriteFile(filePath, []byte(modifiedContent), 0644)
}

func containsAny(s string, substrs []string) bool {
	pattern := strings.Join(substrs, "|")
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}
