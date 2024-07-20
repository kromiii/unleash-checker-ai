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
			if isFlagUsedInFile(path, flags) {
					modifier := modifier.NewModifier(apiKey)
					currentRemovedFlags, err := modifier.ModifyFile(path, flags) // 新しい変数で削除されたフラグを受け取る
					if err != nil {
							return err
					}
					for _, flag := range currentRemovedFlags {
						if !contains(removedFlags, flag) {
								removedFlags = append(removedFlags, flag)
						}
					}
					changedFiles = append(changedFiles, path)
			}
			return nil
	})
	return changedFiles, removedFlags, err
}

func isFlagUsedInFile(path string, flags []string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 1024)
	_, err = file.Read(buf)
	if err != nil {
		return false
	}

	content := string(buf)
	for _, flag := range flags {
		if strings.Contains(content, flag) {
			return true
		}
	}

	return false
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
			if v == str {
					return true
			}
	}
	return false
}
