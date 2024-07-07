package finder

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kromiii/unleash-checker-ai/internal/modifier"
)

func FindAndReplaceFlags(root string, flags []string, apiKey string) ([]string, error) {
	removedFlags := []string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if isFlagUsedInFile(path, flags) {
			modifier := modifier.NewModifier(apiKey)
			removedFlags, err = modifier.ModifyFile(path, flags)
            if err != nil {
                return err
            }
		}
		return nil
	})
	return removedFlags, err
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
