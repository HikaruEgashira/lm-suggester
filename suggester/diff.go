package suggester

import "github.com/sergi/go-diff/diffmatchpatch"

// minimalRangeFromFullAfter は BaseText と LLMAfter 全文の差分から最小置換範囲を返す。
func minimalRangeFromFullAfter(base, after string) (start, end int, afterBlock string, err error) {
	if base == after {
		return 0, 0, "", ErrNoChange
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(base, after, false)
	if len(diffs) == 0 {
		return 0, 0, "", ErrNoChange
	}

	startBase := 0
	startAfter := 0
	i := 0
	for i < len(diffs) && diffs[i].Type == diffmatchpatch.DiffEqual {
		text := diffs[i].Text
		startBase += len(text)
		startAfter += len(text)
		i++
	}
	if i == len(diffs) {
		return 0, 0, "", ErrNoChange
	}

	endBase := len(base)
	endAfter := len(after)
	j := len(diffs) - 1
	for j >= i && diffs[j].Type == diffmatchpatch.DiffEqual {
		text := diffs[j].Text
		endBase -= len(text)
		endAfter -= len(text)
		j--
	}

	if endBase < startBase {
		endBase = startBase
	}
	if endAfter < startAfter {
		endAfter = startAfter
	}

	for endBase < len(base) && endAfter < len(after) && base[endBase] == '\n' && after[endAfter] == '\n' {
		endBase++
		endAfter++
	}

	afterBlock = after[startAfter:endAfter]
	if startBase == endBase && len(afterBlock) == 0 {
		return 0, 0, "", ErrNoChange
	}

	return startBase, endBase, afterBlock, nil
}
