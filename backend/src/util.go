package scsave

import (
	"strings"
)

// 文字列内の連続した空白文字を 1 文字のスペースに変換する
func TrimWordGaps(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, "\u00A0", " ")), " ")
}