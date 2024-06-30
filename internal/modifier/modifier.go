package modifier

import (
	"os"

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

	modifiedContent, err := m.openaiClient.ModifyCode(string(content), unusedFlags)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(modifiedContent), 0644)
}
