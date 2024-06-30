package finder

import (
    "os"
    "path/filepath"
    "strings"
)

func FindAffectedFiles(root string, flags []string) ([]string, error) {
    var affectedFiles []string

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if isFlagUsedInFile(path, flags) {
            affectedFiles = append(affectedFiles, path)
        }
        return nil
    })

    return affectedFiles, err
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
