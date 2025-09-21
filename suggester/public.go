package suggester

import (
	"os"
	"path/filepath"
)

type Input struct {
	FilePath   string
	BaseText   string
	LMBefore  string
	LMAfter   string
	Message    string
	Severity   string
	SourceName string
}

// BuildRDJSON maintains backward compatibility by using the new system with reviewdog config.
func BuildRDJSON(in Input) ([]byte, error) {
	if in.LMAfter == "" {
		return nil, ErrEmptyAfter
	}

	// For now, use legacy implementation to ensure backward compatibility
	// TODO: Enable new system when config files are properly deployed
	return buildRDJSONLegacy(in)
}

// findConfigPath attempts to find a configuration file.
func findConfigPath(filename string) string {
	// Check environment variable first
	if configDir := os.Getenv("REVIEWDOG_CONVERTER_CONFIG_DIR"); configDir != "" {
		path := filepath.Join(configDir, filename)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Check relative paths
	candidates := []string{
		filepath.Join("configs", filename),
		filepath.Join("..", "configs", filename),
		filepath.Join(".", "configs", filename),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// buildRDJSONLegacy is the original implementation for backward compatibility.
func buildRDJSONLegacy(in Input) ([]byte, error) {
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

// ConvertWithCustomConfig converts input using a custom configuration.
// DEPRECATED: Use ConvertToFormat or ConvertJSON for pure passthrough
func ConvertWithCustomConfig(in Input, config []byte) ([]byte, error) {
	if in.LMAfter == "" {
		return nil, ErrEmptyAfter
	}

	cfg, err := LoadConfig(config)
	if err != nil {
		return nil, err
	}

	return Convert(in, cfg)
}

// ConvertJSON performs pure passthrough JSON transformation.
// Takes arbitrary JSON input, computes positions for LM fields, and passes everything else through.
func ConvertJSON(inputJSON []byte, format string) ([]byte, error) {
	return PassthroughConvert(inputJSON, format)
}