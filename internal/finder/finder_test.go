package finder

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FindAndReplaceFlags関数のスタブ
var stubFindAndReplaceFlags func(root string, flags []string, apiKey string) ([]string, []string, error)

func TestFindAndReplaceFlags(t *testing.T) {
	// スタブの設定
	stubFindAndReplaceFlags = func(_ string, _ []string, _ string) ([]string, []string, error) {
		return []string{"file1.txt", "file2.txt"}, []string{"flag1", "flag2"}, nil
	}

	// テスト実行
	changedFiles, removedFlags, err := stubFindAndReplaceFlags("testDir", []string{"flag1", "flag2"}, "dummy-api-key")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, changedFiles, 2)
	assert.Contains(t, changedFiles, "file1.txt")
	assert.Contains(t, changedFiles, "file2.txt")
	assert.ElementsMatch(t, removedFlags, []string{"flag1", "flag2"})
}

func TestIsFlagUsedInFile(t *testing.T) {
	// テストファイルを作成
	tempFile, err := os.CreateTemp("", "test")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("This file contains testFlag"))
	require.NoError(t, err)
	tempFile.Close()

	// テストケース
	testCases := []struct {
		name     string
		flags    []string
		expected bool
	}{
		{"Flag present", []string{"testFlag"}, true},
		{"Flag not present", []string{"nonexistentFlag"}, false},
		{"Multiple flags, one present", []string{"nonexistentFlag", "testFlag"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isFlagUsedInFile(tempFile.Name(), tc.flags)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	assert.True(t, contains(slice, "banana"))
	assert.False(t, contains(slice, "grape"))
}
