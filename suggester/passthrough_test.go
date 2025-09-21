package suggester

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestPassthroughConvert(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		format         string
		checkFields    map[string]interface{}
		checkNotExists []string
		extraFields    map[string]interface{}
	}{
		{
			name: "reviewdog format with passthrough",
			input: `{
				"FilePath": "test.go",
				"BaseText": "func main() {\n\tprintln(\"hello\")\n}\n",
				"LMBefore": "println(\"hello\")",
				"LMAfter": "fmt.Println(\"hello\")",
				"Message": "Use fmt.Println",
				"Severity": "WARNING",
				"SourceName": "test-linter",
				"customField": "custom value",
				"metadata": {
					"version": "1.0",
					"author": "test"
				}
			}`,
			format: "reviewdog",
			checkFields: map[string]interface{}{
				"source.name":                              "test-linter",
				"diagnostics[0].severity":                  "WARNING",
				"diagnostics[0].location.path":             "test.go",
				"diagnostics[0].location.range.start.line": 2.0, // JSON numbers are float64
				"customField":                              "custom value",
			},
			checkNotExists: []string{"FilePath", "BaseText", "LMBefore", "LMAfter"},
			extraFields: map[string]interface{}{
				"metadata.version": "1.0",
				"metadata.author":  "test",
			},
		},
		{
			name: "SARIF format with passthrough",
			input: `{
				"FilePath": "src/app.js",
				"BaseText": "const x = 1;\nconst y = 2;\n",
				"LMBefore": "const y = 2",
				"LMAfter": "let y = 2",
				"Message": "Use let instead of const",
				"Severity": "ERROR",
				"SourceName": "eslint",
				"projectId": "12345",
				"timestamp": "2024-01-01T00:00:00Z"
			}`,
			format: "sarif",
			checkFields: map[string]interface{}{
				"version":                  "2.1.0",
				"runs[0].tool.driver.name": "eslint",
				"runs[0].results[0].level": "error",
				"runs[0].results[0].locations[0].physicalLocation.region.startLine": 2.0,
				"projectId": "12345",
				"timestamp": "2024-01-01T00:00:00Z",
			},
			checkNotExists: []string{"FilePath", "BaseText", "LMBefore", "LMAfter", "Message", "Severity", "SourceName"},
		},
		{
			name: "unknown format - pure passthrough with computed fields",
			input: `{
				"FilePath": "test.py",
				"BaseText": "def foo():\n    pass\n",
				"LMBefore": "    pass",
				"LMAfter": "    return None",
				"arbitrary": "data",
				"nested": {
					"field": "value"
				},
				"array": [1, 2, 3]
			}`,
			format: "custom",
			checkFields: map[string]interface{}{
				"computed.startLine": 2.0,
				"computed.after":     "    return None",
				"arbitrary":          "data",
				"nested.field":       "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PassthroughConvert([]byte(tt.input), tt.format)
			if err != nil {
				t.Fatalf("PassthroughConvert failed: %v", err)
			}

			var output map[string]interface{}
			if err := json.Unmarshal(result, &output); err != nil {
				t.Fatalf("Failed to parse result JSON: %v", err)
			}

			for path, expected := range tt.checkFields {
				actual := getNestedValue(output, path)
				if actual != expected {
					t.Errorf("Field %s: expected %v, got %v", path, expected, actual)
				}
			}

			for _, field := range tt.checkNotExists {
				if _, exists := output[field]; exists {
					t.Errorf("Field %s should not exist at top level after transformation", field)
				}
			}

			for path, expected := range tt.extraFields {
				actual := getNestedValue(output, path)
				if actual != expected {
					t.Errorf("Passthrough field %s: expected %v, got %v", path, expected, actual)
				}
			}
		})
	}
}

func TestPassthroughWithArbitraryJSON(t *testing.T) {
	input := `{
		"FilePath": "test.go",
		"BaseText": "package main\n\nfunc main() {}\n",
		"LMBefore": "func main() {}",
		"LMAfter": "func main() {\n\tfmt.Println(\"Hello\")\n}",
		"customTool": {
			"name": "my-linter",
			"version": "2.0.0",
			"config": {
				"rules": ["no-unused-vars", "semi"],
				"severity": "error"
			}
		},
		"buildInfo": {
			"commit": "abc123",
			"branch": "main",
			"timestamp": 1234567890
		},
		"tags": ["golang", "linting", "ci"],
		"score": 95.5,
		"enabled": true
	}`

	result, err := PassthroughConvert([]byte(input), "reviewdog")
	if err != nil {
		t.Fatalf("PassthroughConvert failed: %v", err)
	}

	var output map[string]interface{}
	if err := json.Unmarshal(result, &output); err != nil {
		t.Fatalf("Failed to parse result JSON: %v", err)
	}

	if getNestedValue(output, "diagnostics[0].location.range.start.line") != 3.0 {
		t.Error("Computed position incorrect")
	}

	passthroughChecks := map[string]interface{}{
		"customTool.name":            "my-linter",
		"customTool.version":         "2.0.0",
		"customTool.config.severity": "error",
		"buildInfo.commit":           "abc123",
		"buildInfo.branch":           "main",
		"buildInfo.timestamp":        1234567890.0,
		"score":                      95.5,
		"enabled":                    true,
	}

	for path, expected := range passthroughChecks {
		actual := getNestedValue(output, path)
		if actual != expected {
			t.Errorf("Passthrough field %s: expected %v (type %T), got %v (type %T)",
				path, expected, expected, actual, actual)
		}
	}

	tags, ok := output["tags"].([]interface{})
	if !ok || len(tags) != 3 {
		t.Error("Tags array not passed through correctly")
	}
}

func TestConvert_DirectAPI(t *testing.T) {
	input := map[string]interface{}{
		"FilePath":   "main.go",
		"BaseText":   "package main\n",
		"LMAfter":    "package main\n\nimport \"fmt\"\n",
		"Message":    "Add import",
		"customData": "preserved",
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	result, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	var output map[string]interface{}
	if err := json.Unmarshal(result, &output); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	if output["customData"] != "preserved" {
		t.Error("Custom data not preserved in passthrough")
	}

	if getNestedValue(output, "diagnostics[0].location.range.start.line") == nil {
		t.Error("Computed fields not added")
	}
}

func getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := interface{}(data)

	for _, part := range parts {
		if idx := strings.Index(part, "["); idx >= 0 {
			fieldName := part[:idx]
			endIdx := strings.Index(part, "]")
			indexStr := part[idx+1 : endIdx]
			index := 0
			if i, err := strconv.Atoi(indexStr); err == nil {
				index = i
			}

			if fieldName != "" {
				if m, ok := current.(map[string]interface{}); ok {
					current = m[fieldName]
				} else {
					return nil
				}
			}

			if arr, ok := current.([]interface{}); ok && index < len(arr) {
				current = arr[index]
			} else {
				return nil
			}
		} else {
			if m, ok := current.(map[string]interface{}); ok {
				current = m[part]
			} else {
				return nil
			}
		}
	}

	return current
}
