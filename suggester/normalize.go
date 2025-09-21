package suggester

import (
	"strings"
	"unicode/utf8"
)

func normalizeText(s string) string {
	if s == "" {
		return s
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func buildLineIndex(s string) []int {
	indexes := []int{0}
	for i, r := range s {
		if r == '\n' {
			indexes = append(indexes, i+1)
		}
	}
	return indexes
}

func offsetToLineCol(s string, offset int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if offset > len(s) {
		offset = len(s)
	}
	lineIdx := buildLineIndex(s)
	lo, hi := 0, len(lineIdx)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if lineIdx[mid] <= offset {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	line := hi + 1
	lineStart := lineIdx[hi]
	col := utf8.RuneCountInString(s[lineStart:offset]) + 1
	return line, col
}