package modifier

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// OpenAIClientMock は OpenAIClientInterface のモック実装です
type OpenAIClientMock struct {
	mock.Mock
}

func (m *OpenAIClientMock) ModifyCode(content string, flags []string) (string, error) {
	args := m.Called(content, flags)
	return args.String(0), args.Error(1)
}

func TestModifier_ModifyFile(t *testing.T) {
	// テスト用の一時ファイルを作成
	content := "line1\nline2\nline3"
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// OpenAIClientMock を設定
	mockClient := new(OpenAIClientMock)
	mockClient.On("ModifyCode", mock.Anything, mock.Anything).Return("2: modified line2", nil)

	// Modifier を作成
	modifier := &Modifier{openaiClient: mockClient}

	// ModifyFile を実行
	err = modifier.ModifyFile(tmpfile.Name(), []string{})
	assert.NoError(t, err)

	// 結果を確認
	modifiedContent, err := os.ReadFile(tmpfile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "line1\nmodified line2\nline3", string(modifiedContent))

	mockClient.AssertExpectations(t)
}

func TestAddLineNumbers(t *testing.T) {
	input := "line1\nline2\nline3"
	expected := "1: line1\n2: line2\n3: line3"
	result := addLineNumbers(input)
	assert.Equal(t, expected, result)
}

func TestStripLineNumbers(t *testing.T) {
	input := "1: line1\n2: line2\n3: line3"
	expected := "line1\nline2\nline3"
	result := stripLineNumbers(input)
	assert.Equal(t, expected, result)
}

func TestApplyDiff(t *testing.T) {
	original := "1: line1\n2: line2\n3: line3"
	diff := "_: new first line\n2: modified line2\n3: \n+: new line4\n+: new last line"
	expected := "new first line\n1: line1\n2: modified line2\nnew line4\nnew last line"
	result, err := applyDiff(original, diff)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
