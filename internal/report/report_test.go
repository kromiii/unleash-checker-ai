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
			name:         "空のフラグリスト",
			staleFlags:   []string{},
			removedFlags: []string{},
			want:         "# Unleash Checker Results Report\n\n## Overview\n**No stale flags were found**\n\n",
		},
		{
			name:         "削除されたフラグあり",
			staleFlags:   []string{"flag1", "flag2"},
			removedFlags: []string{"flag1"},
			want:         "# Unleash Checker Results Report\n\n## Overview\n**Unleash Checker** has detected **2 stale flags**.\n\n## Removed Flags\nThe following flags have been removed from the files (due to being stale):\n\n- `flag1`\n\n## Unfound Flags\nThe following flags were not found in the specified directories:\n\n- `flag2`\n\n",
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
			name:   "完全に異なるスライス",
			sliceA: []string{"a", "b", "c"},
			sliceB: []string{"d", "e", "f"},
			want:   []string{"a", "b", "c"},
		},
		{
			name:   "部分的に重複するスライス",
			sliceA: []string{"a", "b", "c", "d"},
			sliceB: []string{"b", "d"},
			want:   []string{"a", "c"},
		},
		{
			name:   "空のスライス",
			sliceA: []string{},
			sliceB: []string{"a", "b"},
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

// スライスの等価性をチェックするヘルパー関数
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
