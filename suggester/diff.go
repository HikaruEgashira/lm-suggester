package suggester

import "github.com/sergi/go-diff/diffmatchpatch"

func minimalRangeFromFullAfter(base, after string) (start, end int, afterBlock string, err error) {
	dmp := diffmatchpatch.New()
	cp := dmp.DiffCommonPrefix(base, after)
	bs := base[cp:]
	as := after[cp:]
	cs := dmp.DiffCommonSuffix(bs, as)
	start = cp
	end = len(base) - cs
	if end < start {
		end = start
	}
	afterBlock = as[:len(as)-cs]
	if start == end && len(afterBlock) == 0 {
		return 0, 0, "", ErrNoChange
	}
	return start, end, afterBlock, nil
}