package suggester

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
)

// Convert automatically detects whether the input is JSON or JSONL and converts accordingly.
// Takes arbitrary JSON/JSONL input, computes positions for LM fields, and passes everything else through.
func Convert(input []byte, format string) ([]byte, error) {
	if detectJSONL(input) {
		return convertJSONL(input, format)
	}

	result, err := convertJSON(input, format)
	if err == nil {
		return result, nil
	}

	jsonlResult, jsonlErr := convertJSONL(input, format)
	if jsonlErr == nil {
		return jsonlResult, nil
	}

	return nil, err
}

func convertJSON(inputJSON []byte, format string) ([]byte, error) {
	return PassthroughConvert(inputJSON, format)
}

func convertJSONL(inputJSONL []byte, format string) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(inputJSONL))
	var allDiagnostics []interface{}
	var allResults []interface{}
	var sourceName string

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		converted, err := convertJSON(line, format)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		var output map[string]interface{}
		if err := json.Unmarshal(converted, &output); err != nil {
			return nil, fmt.Errorf("line %d: failed to parse converted output: %w", lineNum, err)
		}

		switch format {
		case "reviewdog":
			if diags, ok := output["diagnostics"].([]interface{}); ok {
				allDiagnostics = append(allDiagnostics, diags...)
			}
			if sourceName == "" {
				if source, ok := output["source"].(map[string]interface{}); ok {
					if name, ok := source["name"].(string); ok {
						sourceName = name
					}
				}
			}

		case "sarif":
			if runs, ok := output["runs"].([]interface{}); ok && len(runs) > 0 {
				if run, ok := runs[0].(map[string]interface{}); ok {
					if results, ok := run["results"].([]interface{}); ok {
						allResults = append(allResults, results...)
					}
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
			return nil, fmt.Errorf("JSONL merging not supported for format: %s", format)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read JSONL: %w", err)
	}

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

func detectJSONL(input []byte) bool {
	var singleJSON interface{}
	if err := json.Unmarshal(input, &singleJSON); err == nil {
		return false
	}

	scanner := bufio.NewScanner(bytes.NewReader(input))
	validJSONCount := 0
	totalNonEmptyLines := 0

	for scanner.Scan() {
		line := scanner.Bytes()
		trimmed := bytes.TrimSpace(line)

		if len(trimmed) == 0 {
			continue
		}

		totalNonEmptyLines++

		if len(trimmed) > 0 {
			firstChar := trimmed[0]
			if firstChar != '{' && firstChar != '[' {
				return false
			}
		}

		var lineJSON interface{}
		if err := json.Unmarshal(trimmed, &lineJSON); err == nil {
			validJSONCount++
		}
	}

	return validJSONCount >= 2 && validJSONCount == totalNonEmptyLines
}
