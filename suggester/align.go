package suggester

import (
	"strings"
	"unicode/utf8"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func alignRange(base, before string) (int, int, error) {
	if before == "" {
		return 0, 0, ErrNoMatch
	}
	if idx := strings.Index(base, before); idx >= 0 {
		return idx, idx + len(before), nil
	}

	dmp := diffmatchpatch.New()
	cands := make(map[int]struct{})

	locs := []int{0, len(base) / 4, len(base) / 2, (3 * len(base)) / 4}
	for _, loc := range locs {
		if loc < 0 {
			loc = 0
		}
		if loc > len(base) {
			loc = len(base)
		}
		i := dmp.MatchMain(base, before, loc)
		if i >= 0 {
			cands[i] = struct{}{}
		}
	}

	key := before
	if utf8.RuneCountInString(key) > 64 {
		key = string([]rune(key)[:64])
	}
	if i := strings.Index(base, key); i >= 0 {
		cands[i] = struct{}{}
	}
	if len(before) > 0 {
		runes := []rune(before)
		k2 := runes
		if len(runes) > 64 {
			k2 = runes[len(runes)-64:]
		}
		if i := strings.Index(base, string(k2)); i >= 0 {
			cands[i] = struct{}{}
		}
	}

	longest := longestLine(before)
	if longest != "" {
		if i := strings.Index(base, longest); i >= 0 {
			cands[i] = struct{}{}
		}
	}

	if len(cands) == 0 {
		i := dmp.MatchMain(base, before, 0)
		if i < 0 {
			return 0, 0, ErrNoMatch
		}
		cands[i] = struct{}{}
	}

	bestIdx := -1
	bestDist := -1
	winLen := len(before)
	if winLen > len(base) {
		winLen = len(base)
	}
	for i := range cands {
		start := clampWindowStart(i, winLen, len(base))
		sub := base[start : start+winLen]
		diff := dmp.DiffMain(sub, before, false)
		dist := dmp.DiffLevenshtein(diff)
		if bestIdx < 0 || dist < bestDist || (dist == bestDist && start < bestIdx) {
			bestIdx = start
			bestDist = dist
		}
	}

	if bestIdx < 0 {
		return 0, 0, ErrNoMatch
	}
	return bestIdx, bestIdx + len(before), nil
}

func clampWindowStart(idx, winLen, baseLen int) int {
	start := idx
	if start+winLen > baseLen {
		start = baseLen - winLen
	}
	if start < 0 {
		start = 0
	}
	return start
}

func longestLine(s string) string {
	lines := strings.Split(s, "\n")
	maxLen := 0
	max := ""
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
			max = l
		}
	}
	return max
}