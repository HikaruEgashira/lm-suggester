package suggester

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
)

// ConvertJSON performs pure passthrough JSON transformation.
// Takes arbitrary JSON input, computes positions for LM fields, and passes everything else through.
func ConvertJSON(inputJSON []byte, format string) ([]byte, error) {
	return PassthroughConvert(inputJSON, format)
}

// ConvertJSONL processes JSONL (JSON Lines) format where each line is a separate JSON object.
// It converts each line individually and merges the results into a single output.
func ConvertJSONL(inputJSONL []byte, format string) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(inputJSONL))
	var allDiagnostics []interface{}
	var allResults []interface{}
	var sourceName string

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		// Skip empty lines
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		// Convert each line
		converted, err := ConvertJSON(line, format)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		// Parse the converted output to extract diagnostics/results
		var output map[string]interface{}
		if err := json.Unmarshal(converted, &output); err != nil {
			return nil, fmt.Errorf("line %d: failed to parse converted output: %w", lineNum, err)
		}

		switch format {
		case "reviewdog":
			// Extract diagnostics from this line
			if diags, ok := output["diagnostics"].([]interface{}); ok {
				allDiagnostics = append(allDiagnostics, diags...)
			}
			// Keep the source name from the first line
			if sourceName == "" {
				if source, ok := output["source"].(map[string]interface{}); ok {
					if name, ok := source["name"].(string); ok {
						sourceName = name
					}
				}
			}

		case "sarif":
			// Extract results from runs
			if runs, ok := output["runs"].([]interface{}); ok && len(runs) > 0 {
				if run, ok := runs[0].(map[string]interface{}); ok {
					if results, ok := run["results"].([]interface{}); ok {
						allResults = append(allResults, results...)
					}
					// Keep the tool name from the first line
					if sourceName == "" {
						if tool, ok := run["tool"].(map[string]interface{}); ok {
							if driver, ok := tool["driver"].(map[string]interface{}); ok {
								if name, ok := driver["name"].(string); ok {
									sourceName = name
								}
							}
						}
					}
				}
			}

		default:
			// For unknown formats, merge all computed fields
			return nil, fmt.Errorf("JSONL merging not supported for format: %s", format)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read JSONL: %w", err)
	}

	// Build the final merged output
	var finalOutput map[string]interface{}

	switch format {
	case "reviewdog":
		finalOutput = map[string]interface{}{
			"source": map[string]interface{}{
				"name": sourceName,
			},
			"diagnostics": allDiagnostics,
		}

	case "sarif":
		if sourceName == "" {
			sourceName = "reviewdog-converter"
		}
		finalOutput = map[string]interface{}{
			"version": "2.1.0",
			"$schema": "https://json.schemastore.org/sarif-2.1.0.json",
			"runs": []interface{}{
				map[string]interface{}{
					"tool": map[string]interface{}{
						"driver": map[string]interface{}{
							"name": sourceName,
						},
					},
					"results": allResults,
				},
			},
		}
	}

	return json.Marshal(finalOutput)
}

// DetectJSONL checks if the input appears to be JSONL format.
// Returns true if the input contains multiple JSON objects separated by newlines.
func DetectJSONL(input []byte) bool {
	// Try to parse as single JSON first (including formatted JSON)
	var singleJSON interface{}
	if err := json.Unmarshal(input, &singleJSON); err == nil {
		// It's valid single JSON (formatted or not), not JSONL
		return false
	}

	// If it's not valid single JSON, check if it's JSONL
	scanner := bufio.NewScanner(bytes.NewReader(input))
	validJSONCount := 0
	totalNonEmptyLines := 0

	for scanner.Scan() {
		line := scanner.Bytes()
		trimmed := bytes.TrimSpace(line)

		// Skip empty lines
		if len(trimmed) == 0 {
			continue
		}

		totalNonEmptyLines++

		// Check if line starts with typical JSON tokens
		if len(trimmed) > 0 {
			firstChar := trimmed[0]
			// If line doesn't start with '{' or '[', it's likely part of formatted JSON
			if firstChar != '{' && firstChar != '[' {
				// Could be part of a formatted JSON, not JSONL
				return false
			}
		}

		// Try to parse this line as JSON
		var lineJSON interface{}
		if err := json.Unmarshal(trimmed, &lineJSON); err == nil {
			validJSONCount++
		}
	}

	// It's JSONL if we found at least 2 valid JSON objects and all non-empty lines are valid JSON
	return validJSONCount >= 2 && validJSONCount == totalNonEmptyLines
}

// ConvertAuto automatically detects whether the input is JSON or JSONL and converts accordingly.
func ConvertAuto(input []byte, format string) ([]byte, error) {
	// First check if it's JSONL
	if DetectJSONL(input) {
		return ConvertJSONL(input, format)
	}

	// Try parsing as single JSON
	result, err := ConvertJSON(input, format)
	if err == nil {
		return result, nil
	}

	// If single JSON failed, try JSONL as fallback
	// This handles cases where detection might miss edge cases
	jsonlResult, jsonlErr := ConvertJSONL(input, format)
	if jsonlErr == nil {
		return jsonlResult, nil
	}

	// Return the original error from single JSON attempt
	return nil, err
}