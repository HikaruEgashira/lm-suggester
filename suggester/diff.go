package suggester

import "github.com/sergi/go-diff/diffmatchpatch"

func minimalRangeFromFullAfter(base, after string) (start, end int, afterBlock string, err error) {
	dmp := diffmatchpatch.New()

	// DiffCommonPrefix returns rune count, need to convert to byte offset
	cpRunes := dmp.DiffCommonPrefix(base, after)
	baseRunes := []rune(base)

	// Convert rune index to byte offset for prefix
	cpBytes := 0
	if cpRunes > 0 && cpRunes <= len(baseRunes) {
		cpBytes = len(string(baseRunes[:cpRunes]))
	}

	// Get the remaining parts after common prefix
	bs := base[cpBytes:]
	as := after[cpBytes:]

	// DiffCommonSuffix also returns rune count
	csRunes := dmp.DiffCommonSuffix(bs, as)

	// Convert rune count to byte offset for suffix
	csBytes := 0
	if csRunes > 0 {
		bsRunes := []rune(bs)
		if csRunes <= len(bsRunes) {
			csBytes = len(string(bsRunes[len(bsRunes)-csRunes:]))
		}
	}

	start = cpBytes
	end = len(base) - csBytes
	if end < start {
		end = start
	}

	// Calculate afterBlock using byte offsets
	afterBlockEndBytes := len(as) - csBytes
	if afterBlockEndBytes < 0 {
		afterBlockEndBytes = 0
	}
	afterBlock = as[:afterBlockEndBytes]

	if start == end && len(afterBlock) == 0 {
		return 0, 0, "", ErrNoChange
	}
	return start, end, afterBlock, nil
}