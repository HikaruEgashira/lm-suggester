package suggester

// CoreResult represents the minimal transformation result containing
// only the essential computed values (position calculations).
// All other input fields are passed through unchanged.
type CoreResult struct {
	FilePath    string
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Before      string // Actual text before change
	After       string // Actual text after change
}

// ExtractCore performs core processing (diff detection and position calculation)
// without any format-specific transformations.
func ExtractCore(filePath, baseText, lmBefore, lmAfter string) (*CoreResult, error) {
	// Normalize line endings
	baseText = normalizeText(baseText)
	lmAfter = normalizeText(lmAfter)
	if lmBefore != "" {
		lmBefore = normalizeText(lmBefore)
	}

	var startOffset, endOffset int
	var beforeText, afterText string

	if lmBefore != "" {
		// If LMBefore is provided, align it with the base text
		start, end, err := alignRange(baseText, lmBefore)
		if err != nil {
			return nil, err
		}
		startOffset = start
		endOffset = end
		beforeText = lmBefore
		afterText = lmAfter
	} else {
		// If LMBefore is not provided, find minimal diff range
		start, end, afterBlock, err := minimalRangeFromFullAfter(baseText, lmAfter)
		if err != nil {
			return nil, err
		}
		startOffset = start
		endOffset = end
		beforeText = baseText[startOffset:endOffset]
		afterText = afterBlock
	}

	// Convert byte offsets to line/column
	startLine, startCol := offsetToLineCol(baseText, startOffset)
	endLine, endCol := offsetToLineCol(baseText, endOffset)

	return &CoreResult{
		FilePath:    filePath,
		StartLine:   startLine,
		StartColumn: startCol,
		EndLine:     endLine,
		EndColumn:   endCol,
		Before:      beforeText,
		After:       afterText,
	}, nil
}