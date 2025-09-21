package suggester

import "github.com/sergi/go-diff/diffmatchpatch"

func minimalRangeFromFullAfter(base, after string) (start, end int, afterBlock string, err error) {
	dmp := diffmatchpatch.New()

	cpRunes := dmp.DiffCommonPrefix(base, after)
	baseRunes := []rune(base)

	cpBytes := 0
	if cpRunes > 0 && cpRunes <= len(baseRunes) {
		cpBytes = len(string(baseRunes[:cpRunes]))
	}

	bs := base[cpBytes:]
	as := after[cpBytes:]

	csRunes := dmp.DiffCommonSuffix(bs, as)

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
