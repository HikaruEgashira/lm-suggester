package suggester

type Input struct {
	FilePath   string
	BaseText   string
	LMBefore  string
	LMAfter   string
	Message    string
	Severity   string
	SourceName string
}

// BuildRDJSON builds reviewdog JSON format output.
func BuildRDJSON(in Input) ([]byte, error) {
	if in.LMAfter == "" {
		return nil, ErrEmptyAfter
	}
	base := normalizeText(in.BaseText)
	before := normalizeText(in.LMBefore)
	after := normalizeText(in.LMAfter)

	var (
		start, end int
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

// ConvertToFormat converts input to a specific format using pure passthrough.
// Supported formats: "reviewdog", "sarif", and any custom format (passthrough with computed fields)
func ConvertToFormat(in Input, format string) ([]byte, error) {
	if in.LMAfter == "" {
		return nil, ErrEmptyAfter
	}

	// Convert Input to JSON for passthrough processing
	inputJSON, err := ConvertInputToJSON(in)
	if err != nil {
		return nil, err
	}

	return PassthroughConvert(inputJSON, format)
}

// ConvertJSON performs pure passthrough JSON transformation.
// Takes arbitrary JSON input, computes positions for LM fields, and passes everything else through.
func ConvertJSON(inputJSON []byte, format string) ([]byte, error) {
	return PassthroughConvert(inputJSON, format)
}