package report

import (
	"fmt"
	"strings"
)

func CreateSummary(staleFlags []string, removedFlags []string) string {
	var output strings.Builder

	output.WriteString("# Unleash Checker の結果報告\n\n")
	
	output.WriteString("## 概要\n")
	if len(staleFlags) == 0 {
		output.WriteString("**古いフラグは見つかりませんでした**\n\n")
		return output.String()
	}
	output.WriteString(fmt.Sprintf("**Unleash Checker** が **%d個の古いフラグ** を検出しました。\n\n", len(staleFlags)))

	if len(removedFlags) > 0 {
		output.WriteString("## 削除されたフラグ\n")
		output.WriteString("以下のフラグがファイルから削除されました（古くなっているため）:\n\n")
		for _, flag := range removedFlags {
			output.WriteString(fmt.Sprintf("- `%s`\n", flag))
		}
		output.WriteString("\n")
	}

	unfoundFlags := difference(staleFlags, removedFlags)
	if len(unfoundFlags) > 0 {
		output.WriteString("## 見つからなかったフラグ\n")
		output.WriteString("以下のフラグは指定されたディレクトリで見つかりませんでした:\n\n")
		for _, flag := range unfoundFlags {
			output.WriteString(fmt.Sprintf("- `%s`\n", flag))
		}
		output.WriteString("\n")
	}

	output.WriteString("## アクション項目\n")
	output.WriteString("1. 変更内容を確認してください\n")
	output.WriteString("2. 問題がなければ、リポジトリにコミットしてください\n\n")

	output.WriteString("## 注意事項\n")
	output.WriteString("これらのフラグをまだ使用したい場合は、フラグタイプの変更を検討してください。\n")
	output.WriteString("詳細については以下のドキュメントを参照してください:\n")
	output.WriteString("https://docs.getunleash.io/reference/technical-debt\n\n")

	output.WriteString("---\n\n")
	output.WriteString("このPRは **Unleash Checker AI GitHub Action** によって自動生成されました。\n")

	return output.String()
}

func difference(sliceA, sliceB []string) []string {
	var diff []string

	for _, a := range sliceA {
		found := false
		for _, b := range sliceB {
			if a == b {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, a)
		}
	}

	return diff
}
