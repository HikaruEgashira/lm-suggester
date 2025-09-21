package suggester

import (
	"sort"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// alignRange は BaseText 中で LLMBefore が最も適合する範囲 [start,end) を返す。
func alignRange(base, before string) (int, int, error) {
	if before == "" {
		return 0, 0, ErrNoMatch
	}
	if len(base) == 0 {
		return 0, 0, ErrNoMatch
	}

	if idx := strings.Index(base, before); idx >= 0 {
		end := idx + len(before)
		if end > len(base) {
			end = len(base)
		}
		return idx, end, nil
	}

	dmp := diffmatchpatch.New()
	candidates := collectCandidates(base, before, dmp)
	if len(candidates) == 0 {
		pos := dmp.MatchMain(base, before, 0)
		if pos < 0 {
			return 0, 0, ErrNoMatch
		}
		candidates = append(candidates, pos)
	}

	bestStart := -1
	bestScore := -1
	bestLen := 0
	window := len(before)
	if window == 0 {
		return 0, 0, ErrNoMatch
	}
	for _, cand := range candidates {
		for _, shift := range []int{-128, -64, -32, -8, 0, 8, 32, 64, 128} {
			start := cand + shift
			if start < 0 {
				start = 0
			}
			if start > len(base) {
				start = len(base)
			}
			end := start + window
			if end > len(base) {
				end = len(base)
			}
			if end < start {
				continue
			}
			segment := base[start:end]
			if len(segment) == 0 {
				continue
			}
			diffs := dmp.DiffMain(segment, before, false)
			score := dmp.DiffLevenshtein(diffs)
			if bestStart < 0 || score < bestScore || (score == bestScore && start < bestStart) {
				bestStart = start
				bestScore = score
				bestLen = len(segment)
			}
		}
	}

	if bestStart < 0 {
		return 0, 0, ErrNoMatch
	}

	end := bestStart + bestLen
	if end > len(base) {
		end = len(base)
	}
	return bestStart, end, nil
}

func collectCandidates(base, before string, dmp *diffmatchpatch.DiffMatchPatch) []int {
	candidateSet := make(map[int]struct{})
	lengths := []int{0, len(base) / 4, len(base) / 2, (3 * len(base)) / 4}
	for _, loc := range lengths {
		if loc < 0 {
			loc = 0
		}
		if loc > len(base) {
			loc = len(base)
		}
		idx := dmp.MatchMain(base, before, loc)
		if idx >= 0 {
			candidateSet[idx] = struct{}{}
		}
	}

	addOccurrences := func(fragment string) {
		if fragment == "" {
			return
		}
		offset := 0
		for offset <= len(base) {
			i := strings.Index(base[offset:], fragment)
			if i < 0 {
				break
			}
			pos := offset + i
			candidateSet[pos] = struct{}{}
			offset = pos + 1
		}
	}

	runes := []rune(before)
	if len(runes) > 0 {
		addOccurrences(string(runes[:minInt(len(runes), 64)]))
		addOccurrences(string(runes[maxInt(0, len(runes)-64):]))
	}
	addOccurrences(longestLine(before))

	res := make([]int, 0, len(candidateSet))
	for idx := range candidateSet {
		res = append(res, idx)
	}
	sort.Ints(res)
	return res
}

func longestLine(s string) string {
	lines := strings.Split(s, "\n")
	longest := ""
	maxLen := -1
	for _, line := range lines {
		l := len(line)
		if l > maxLen {
			maxLen = l
			longest = line
		}
	}
	return longest
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
