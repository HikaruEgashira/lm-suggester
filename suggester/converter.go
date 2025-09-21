package suggester

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Convert performs generic JSON transformation using the specified configuration.
func Convert(input Input, config *TransformConfig) ([]byte, error) {
	// Extract core results (position calculations)
	core, err := ExtractCore(input.FilePath, input.BaseText, input.LMBefore, input.LMAfter)
	if err != nil {
		return nil, fmt.Errorf("failed to extract core: %w", err)
	}

	// Add byte offsets to core (useful for some formats)
	extendedCore := extendCoreWithOffsets(core, input.BaseText)

	// Convert input to map for pass-through
	inputMap, err := structToMap(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to map: %w", err)
	}

	// Create transformation context
	transform := &Transform{
		Input:  inputMap,
		Core:   extendedCore,
		Config: config,
	}

	// Execute transformation
	return transform.Execute()
}

// ConvertWithConfig loads a configuration file and performs the transformation.
func ConvertWithConfig(input Input, configPath string) ([]byte, error) {
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config, err := LoadConfig(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return Convert(input, config)
}

// ConvertWithFormat converts using a predefined format configuration.
func ConvertWithFormat(input Input, format string) ([]byte, error) {
	// Determine config path based on format
	configDir := getConfigDir()
	configPath := filepath.Join(configDir, format+".json")

	return ConvertWithConfig(input, configPath)
}

// getConfigDir returns the configuration directory path.
func getConfigDir() string {
	// Try to find the configs directory relative to the current executable
	// or use a default location
	if dir := os.Getenv("REVIEWDOG_CONVERTER_CONFIG_DIR"); dir != "" {
		return dir
	}

	// Default to configs directory in the module root
	return "configs"
}

// structToMap converts a struct to a map for generic processing.
func structToMap(v interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Use actual field name for mapping
		result[field.Name] = fieldValue.Interface()
	}

	return result, nil
}

// extendCoreWithOffsets adds byte offset information to CoreResult.
type ExtendedCore struct {
	*CoreResult
	StartOffset int // Byte offset of start position
	EndOffset   int // Byte offset of end position
}

func extendCoreWithOffsets(core *CoreResult, baseText string) *CoreResult {
	// Calculate byte offsets from line/column
	// This is a simplified implementation
	// In the future, we could extend CoreResult to include these offsets
	// for formats that need byte positions instead of line/column
	_ = lineColToOffset(baseText, core.StartLine, core.StartColumn)
	_ = lineColToOffset(baseText, core.EndLine, core.EndColumn)

	// For now, just return the core as-is
	// This function is a placeholder for future extensions
	return core
}

// lineColToOffset converts line/column to byte offset.
func lineColToOffset(text string, line, col int) int {
	if line <= 0 || col <= 0 {
		return 0
	}

	offset := 0
	currentLine := 1
	currentCol := 1

	for i, r := range text {
		if currentLine == line && currentCol == col {
			return i
		}

		if r == '\n' {
			currentLine++
			currentCol = 1
		} else {
			currentCol++
		}
		offset = i
	}

	return offset
}