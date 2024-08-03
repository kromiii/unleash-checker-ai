package finder

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kromiii/unleash-checker-ai/internal/modifier"
)

func FindAndReplaceFlags(root string, flags []string, apiKey string) ([]string, []string, error) {
	removedFlags := []string{}
	changedFiles := []string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		foundFlags := isFlagUsedInFile(path, flags)
		if len(foundFlags) > 0 {
			modifier := modifier.NewModifier(apiKey)
			err := modifier.ModifyFile(path, foundFlags) // 新しい変数で削除されたフラグを受け取る
			if err != nil {
				return err
			}
			removedFlags = append(removedFlags, foundFlags...)
			changedFiles = append(changedFiles, path)
		}
		return nil
	})
	removedFlags = removeDuplicateStrings(removedFlags)
	return changedFiles, removedFlags, err
}

func isFlagUsedInFile(path string, flags []string) []string {
	foundFlags := []string{}
	content, err := os.ReadFile(path) // ファイル全体を読み込む
	if err != nil {
		return foundFlags
	}

	for _, flag := range flags {
		if strings.Contains(string(content), flag) {
			foundFlags = append(foundFlags, flag)
		}
	}

	return foundFlags
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func removeDuplicateStrings(s []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, str := range s {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}
	return result
}
