package modifier

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) ModifyCode(content, instruction string) (string, error) {
	args := m.Called(content, instruction)
	return args.String(0), args.Error(1)
}

var _ OpenAIClientInterface = (*MockOpenAIClient)(nil)

func TestModifyFile(t *testing.T) {
	// テスト用の一時ファイルを作成
	tempFile, err := os.CreateTemp("", "test_file_*.go")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// テスト用のコンテンツを書き込む
	testContent := `package main

func main() {
	if featureFlag1 {
		// some code
	}
	if featureFlag2 {
		// some other code
	}
}
`
	err = os.WriteFile(tempFile.Name(), []byte(testContent), 0644)
	assert.NoError(t, err)

	// モックOpenAIクライアントを作成
	mockClient := new(MockOpenAIClient)
	mockClient.On("ModifyCode", mock.Anything, mock.Anything).Return(testContent, nil)

	// Modifierを作成
	m := &Modifier{openaiClient: mockClient}

	// テスト実行
	unusedFlags := []string{"featureFlag1", "featureFlag2"}
	removedFlags, err := m.ModifyFile(tempFile.Name(), unusedFlags)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, unusedFlags, removedFlags)

	// ファイルの内容を確認
	modifiedContent, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)
	expectedContent := `package main

// This feature flag is stale and can be removed: featureFlag1
func main() {
	if featureFlag1 {
		// some code
	}
// This feature flag is stale and can be removed: featureFlag2
	if featureFlag2 {
		// some other code
	}
}
`
	assert.Equal(t, expectedContent, string(modifiedContent))

	// モックの呼び出しを確認
	mockClient.AssertExpectations(t)
}

func TestFindMatchedFlag(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		flags    []string
		expected bool
		matched  string
	}{
		{
			name:     "Match found",
			input:    "if featureFlag1 {",
			flags:    []string{"featureFlag1", "featureFlag2"},
			expected: true,
			matched:  "featureFlag1",
		},
		{
			name:     "No match found",
			input:    "if someOtherFlag {",
			flags:    []string{"featureFlag1", "featureFlag2"},
			expected: false,
			matched:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, matched := findMatchedFlag(tc.input, tc.flags)
			assert.Equal(t, tc.expected, found)
			assert.Equal(t, tc.matched, matched)
		})
	}
}