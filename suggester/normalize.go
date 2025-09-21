package suggester

import (
	"strings"
	"unicode/utf8"
)

// normalizeText は改行コードを LF に統一する。
func normalizeText(s string) string {
	if s == "" {
		return s
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

// buildLineIndex は各行の先頭バイトオフセットを返す。
func buildLineIndex(s string) []int {
	idx := []int{0}
	for i, r := range s {
		if r == '\n' {
			idx = append(idx, i+1)
		}
	}
	return idx
}

// offsetToLineCol は 0-origin バイトオフセットを 1-origin の行・列に変換する。
func offsetToLineCol(s string, offset int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if offset > len(s) {
		offset = len(s)
	}
	idx := buildLineIndex(s)
	lineIdx := 0
	for lineIdx+1 < len(idx) && idx[lineIdx+1] <= offset {
		lineIdx++
	}
	line := lineIdx + 1
	lineStart := idx[lineIdx]
	col := utf8.RuneCountInString(s[lineStart:offset]) + 1
	return line, col
}
