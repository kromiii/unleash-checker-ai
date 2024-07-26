package report

import (
	"strings"
	"testing"
)

func TestCreateSummary(t *testing.T) {
	tests := []struct {
		name         string
		staleFlags   []string
		removedFlags []string
		want         string
	}{
		{
			name:         "古いフラグがない場合",
			staleFlags:   []string{},
			removedFlags: []string{},
			want:         "# Unleash Checker の結果報告\n\n## 概要\n**古いフラグは見つかりませんでした**\n\n",
		},
		{
			name:         "古いフラグがあり、一部が削除された場合",
			staleFlags:   []string{"flag1", "flag2", "flag3"},
			removedFlags: []string{"flag1", "flag2"},
			want: "# Unleash Checker の結果報告\n\n## 概要\n**Unleash Checker** が **3個の古いフラグ** を検出しました。\n\n" +
				"## 削除されたフラグ\n以下のフラグがファイルから削除されました（古くなっているため）:\n\n- `flag1`\n- `flag2`\n\n" +
				"## 見つからなかったフラグ\n以下のフラグは指定されたディレクトリで見つかりませんでした:\n\n- `flag3`\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateSummary(tt.staleFlags, tt.removedFlags)
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("CreateSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name   string
		sliceA []string
		sliceB []string
		want   []string
	}{
		{
			name:   "スライスAにのみ存在する要素がある場合",
			sliceA: []string{"a", "b", "c"},
			sliceB: []string{"b", "c"},
			want:   []string{"a"},
		},
		{
			name:   "両方のスライスが同じ場合",
			sliceA: []string{"a", "b", "c"},
			sliceB: []string{"a", "b", "c"},
			want:   []string{},
		},
		{
			name:   "スライスAが空の場合",
			sliceA: []string{},
			sliceB: []string{"a", "b", "c"},
			want:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := difference(tt.sliceA, tt.sliceB)
			if !equalSlices(got, tt.want) {
				t.Errorf("difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
