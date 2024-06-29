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
        if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".py") {
            // ファイルの内容を読み込み、フラグの使用を確認
            // この実装は簡略化しています
            affectedFiles = append(affectedFiles, path)
        }
        return nil
    })

    return affectedFiles, err
}
