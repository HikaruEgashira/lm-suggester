package suggester

import (
	"encoding/json"
	"fmt"
	"os"
)

// PassthroughTransformer provides pure passthrough JSON transformation
// Only computes necessary fields (position calculations) and passes everything else through
type PassthroughTransformer struct {
	// Function to compute position and diff information
	ComputeCore func(input map[string]interface{}) (*CoreResult, error)

	// Function to inject computed values into the output
	InjectComputed func(output map[string]interface{}, core *CoreResult) error
}

// Transform performs passthrough transformation with minimal computation
func (pt *PassthroughTransformer) Transform(input []byte) ([]byte, error) {
	// Parse input JSON to map
	var inputMap map[string]interface{}
	if err := json.Unmarshal(input, &inputMap); err != nil {
		return nil, fmt.Errorf("failed to parse input JSON: %w", err)
	}

	// Compute core values (position calculations)
	core, err := pt.ComputeCore(inputMap)
	if err != nil {
		return nil, fmt.Errorf("failed to compute core values: %w", err)
	}

	// Start with input as base (complete passthrough)
	output := inputMap

	// Inject only the computed values
	if err := pt.InjectComputed(output, core); err != nil {
		return nil, fmt.Errorf("failed to inject computed values: %w", err)
	}

	// Marshal back to JSON
	return json.Marshal(output)
}

// StandardComputeCore extracts necessary fields and computes positions
func StandardComputeCore(input map[string]interface{}) (*CoreResult, error) {
	// Extract required fields - support both camelCase and snake_case
	filePath := getFieldWithFallback(input, "FilePath", "file_path")
	baseText := getFieldWithFallback(input, "BaseText", "base_text")
	lmBefore := getFieldWithFallback(input, "LMBefore", "lm_before")
	lmAfter := getFieldWithFallback(input, "LMAfter", "lm_after")

	if lmAfter == "" {
		return nil, fmt.Errorf("LMAfter/lm_after is required")
	}

	// BaseText is required for proper diff calculation
	if baseText == "" {
		// Try to read from file if FilePath is provided
		if filePath != "" {
			content, err := os.ReadFile(filePath)
			if err == nil {
				baseText = string(content)
			} else {
				return nil, fmt.Errorf("BaseText/base_text is required (failed to read from %s: %w)", filePath, err)
			}
		} else {
			return nil, fmt.Errorf("BaseText/base_text is required")
		}
	}

	// Compute positions using existing logic
	return ExtractCore(filePath, baseText, lmBefore, lmAfter)
}

// getFieldWithFallback tries to get a field with multiple possible names
func getFieldWithFallback(input map[string]interface{}, names ...string) string {
	for _, name := range names {
		if val, ok := input[name].(string); ok && val != "" {
			return val
		}
	}
	return ""
}

// ReviewdogInjector injects computed values for reviewdog format
func ReviewdogInjector(output map[string]interface{}, core *CoreResult) error {
	// Create diagnostics structure if not exists
	if _, exists := output["diagnostics"]; !exists {
		output["diagnostics"] = []interface{}{}
	}

	diagnostics, ok := output["diagnostics"].([]interface{})
	if !ok {
		diagnostics = []interface{}{}
	}

	// Create diagnostic entry with computed positions
	diagnostic := map[string]interface{}{
		"location": map[string]interface{}{
			"path": core.FilePath,
			"range": map[string]interface{}{
				"start": map[string]interface{}{
					"line":   core.StartLine,
					"column": core.StartColumn,
				},
				"end": map[string]interface{}{
					"line":   core.EndLine,
					"column": core.EndColumn,
				},
			},
		},
	}

	// Copy other fields from input to diagnostic
	afterText := core.After
	// Remove trailing newline if present for message formatting
	if len(afterText) > 0 && afterText[len(afterText)-1] == '\n' {
		afterText = afterText[:len(afterText)-1]
	}

	// Check both Message (capital) and message (lowercase) for compatibility
	message := getFieldWithFallback(output, "Message", "message")
	if message != "" {
		diagnostic["message"] = fmt.Sprintf("%v\n```suggestion\n%s\n```", message, afterText)
		delete(output, "Message") // Remove from top level
		delete(output, "message") // Remove lowercase version too
	} else {
		// Default message with suggestion
		diagnostic["message"] = fmt.Sprintf("Replace code with suggestion\n```suggestion\n%s\n```", afterText)
	}

	if severity, exists := output["Severity"]; exists {
		diagnostic["severity"] = severity
		delete(output, "Severity") // Remove from top level
	}

	// Add diagnostic
	if len(diagnostics) == 0 {
		diagnostics = append(diagnostics, diagnostic)
	} else {
		// Update first diagnostic with computed values
		if first, ok := diagnostics[0].(map[string]interface{}); ok {
			// Merge computed values into existing diagnostic
			if loc, ok := diagnostic["location"].(map[string]interface{}); ok {
				first["location"] = loc
			}
			if msg, ok := diagnostic["message"].(string); ok {
				first["message"] = msg
			}
			if sev, ok := diagnostic["severity"]; ok {
				first["severity"] = sev
			}
		}
	}

	output["diagnostics"] = diagnostics

	// Ensure source exists
	if _, exists := output["source"]; !exists {
		output["source"] = map[string]interface{}{}
	}
	if source, ok := output["source"].(map[string]interface{}); ok {
		if sourceName, exists := output["SourceName"]; exists {
			source["name"] = sourceName
			delete(output, "SourceName") // Remove from top level
		} else if _, exists := source["name"]; !exists {
			source["name"] = "reviewdog-converter"
		}
	}

	// Clean up consumed fields
	delete(output, "FilePath")
	delete(output, "BaseText")
	delete(output, "LMBefore")
	delete(output, "LMAfter")

	return nil
}

// SARIFInjector injects computed values for SARIF format
func SARIFInjector(output map[string]interface{}, core *CoreResult) error {
	// Ensure SARIF structure exists
	if _, exists := output["version"]; !exists {
		output["version"] = "2.1.0"
	}
	if _, exists := output["$schema"]; !exists {
		output["$schema"] = "https://json.schemastore.org/sarif-2.1.0.json"
	}
	if _, exists := output["runs"]; !exists {
		output["runs"] = []interface{}{}
	}

	runs, ok := output["runs"].([]interface{})
	if !ok || len(runs) == 0 {
		runs = []interface{}{
			map[string]interface{}{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name": "reviewdog-converter",
					},
				},
				"results": []interface{}{},
			},
		}
	}

	run := runs[0].(map[string]interface{})

	// Update tool name if SourceName exists
	if sourceName, exists := output["SourceName"]; exists {
		if tool, ok := run["tool"].(map[string]interface{}); ok {
			if driver, ok := tool["driver"].(map[string]interface{}); ok {
				driver["name"] = sourceName
			}
		}
		delete(output, "SourceName")
	}

	// Create or update results
	if _, exists := run["results"]; !exists {
		run["results"] = []interface{}{}
	}

	results := run["results"].([]interface{})

	// Check both Message (capital) and message (lowercase)
	messageText := getFieldWithFallback(output, "Message", "message")
	result := map[string]interface{}{
		"message": map[string]interface{}{
			"text": messageText,
		},
		"locations": []interface{}{
			map[string]interface{}{
				"physicalLocation": map[string]interface{}{
					"artifactLocation": map[string]interface{}{
						"uri": core.FilePath,
					},
					"region": map[string]interface{}{
						"startLine":   core.StartLine,
						"startColumn": core.StartColumn,
						"endLine":     core.EndLine,
						"endColumn":   core.EndColumn,
					},
				},
			},
		},
	}

	// Convert severity
	if severity, exists := output["Severity"]; exists {
		level := "warning"
		switch severity {
		case "ERROR":
			level = "error"
		case "WARNING":
			level = "warning"
		case "INFO":
			level = "note"
		}
		result["level"] = level
		delete(output, "Severity")
	}

	// Add fix suggestion if present
	if core.After != "" {
		result["fixes"] = []interface{}{
			map[string]interface{}{
				"description": map[string]interface{}{
					"text": "Apply suggested fix",
				},
				"artifactChanges": []interface{}{
					map[string]interface{}{
						"artifactLocation": map[string]interface{}{
							"uri": core.FilePath,
						},
						"replacements": []interface{}{
							map[string]interface{}{
								"deletedRegion": map[string]interface{}{
									"startLine":   core.StartLine,
									"startColumn": core.StartColumn,
									"endLine":     core.EndLine,
									"endColumn":   core.EndColumn,
								},
								"insertedContent": map[string]interface{}{
									"text": core.After,
								},
							},
						},
					},
				},
			},
		}
	}

	if len(results) == 0 {
		results = append(results, result)
	} else {
		// Update first result
		results[0] = result
	}

	run["results"] = results
	output["runs"] = runs

	// Clean up consumed fields
	delete(output, "FilePath")
	delete(output, "BaseText")
	delete(output, "LMBefore")
	delete(output, "LMAfter")
	delete(output, "Message")
	delete(output, "message") // Also remove lowercase version

	return nil
}

// PassthroughConvert converts JSON using pure passthrough with minimal computation
func PassthroughConvert(input []byte, format string) ([]byte, error) {
	var injector func(map[string]interface{}, *CoreResult) error

	switch format {
	case "reviewdog":
		injector = ReviewdogInjector
	case "sarif":
		injector = SARIFInjector
	default:
		// For unknown formats, just pass through with computed fields added
		injector = func(output map[string]interface{}, core *CoreResult) error {
			// Add computed positions as-is
			output["computed"] = map[string]interface{}{
				"startLine":   core.StartLine,
				"startColumn": core.StartColumn,
				"endLine":     core.EndLine,
				"endColumn":   core.EndColumn,
				"before":      core.Before,
				"after":       core.After,
			}
			return nil
		}
	}

	transformer := &PassthroughTransformer{
		ComputeCore:    StandardComputeCore,
		InjectComputed: injector,
	}

	return transformer.Transform(input)
}