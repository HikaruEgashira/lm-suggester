package suggester

// Input は BuildRDJSON への入力を表す。
type Input struct {
	FilePath   string
	BaseText   string
	LLMBefore  string
	LLMAfter   string
	Message    string
	Severity   string
	SourceName string
}

// BuildRDJSON は入力に基づき reviewdog の RDJSON フォーマットを生成する。
func BuildRDJSON(in Input) ([]byte, error) {
	if in.LLMAfter == "" {
		return nil, ErrEmptyAfter
	}

	base := normalizeText(in.BaseText)
	before := normalizeText(in.LLMBefore)
	after := normalizeText(in.LLMAfter)

	var (
		start      int
		end        int
		afterBlock string
		err        error
	)

	if before != "" {
		start, end, err = alignRange(base, before)
		if err != nil {
			return nil, err
		}
		afterBlock = after
	} else {
		start, end, afterBlock, err = minimalRangeFromFullAfter(base, after)
		if err != nil {
			return nil, err
		}
	}

	startLine, startCol := offsetToLineCol(base, start)
	endLine, endCol := offsetToLineCol(base, end)

	msg := makeMessage(in.Message, afterBlock)
	return marshalRDJSON(in.SourceName, in.FilePath, msg, startLine, startCol, endLine, endCol, in.Severity)
}
