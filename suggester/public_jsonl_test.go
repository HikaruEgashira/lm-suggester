package suggester

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDetectJSONL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "single JSON object",
			input:    `{"key": "value"}`,
			expected: false,
		},
		{
			name: "formatted JSON",
			input: `{
  "key": "value",
  "nested": {
    "field": "data"
  }
}`,
			expected: false,
		},
		{
			name: "JSONL with two lines",
			input: `{"key": "value1"}
{"key": "value2"}`,
			expected: true,
		},
		{
			name: "JSONL with empty lines",
			input: `{"key": "value1"}

{"key": "value2"}

`,
			expected: true,
		},
		{
			name:     "invalid JSON",
			input:    `not json at all`,
			expected: false,
		},
		{
			name: "mixed valid and invalid lines",
			input: `{"key": "value1"}
not json
{"key": "value2"}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectJSONL([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("detectJSONL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertJSONL(t *testing.T) {
	jsonl := `{"FilePath": "main.go", "BaseText": "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n", "LMBefore": "\tprintln(\"hello\")", "LMAfter": "\tfmt.Println(\"Hello\")", "Message": "Use fmt.Println"}
{"FilePath": "test.go", "BaseText": "package main\n\nfunc test() {\n\treturn\n}\n", "LMBefore": "func test()", "LMAfter": "func Test()", "Message": "Export function"}`

	result, err := convertJSONL([]byte(jsonl), "reviewdog")
	if err != nil {
		t.Fatalf("convertJSONL() error = %v", err)
	}

	var output map[string]interface{}
	if err := json.Unmarshal(result, &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	diagnostics, ok := output["diagnostics"].([]interface{})
	if !ok {
		t.Fatal("No diagnostics in output")
	}

	if len(diagnostics) != 2 {
		t.Errorf("Expected 2 diagnostics, got %d", len(diagnostics))
	}

	if _, ok := output["source"]; !ok {
		t.Error("No source in output")
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "single JSON",
			input: `{"FilePath": "main.go", "BaseText": "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n", "LMAfter": "\tfmt.Println(\"Hello\")"}`,
		},
		{
			name: "formatted JSON",
			input: `{
  "FilePath": "main.go",
  "BaseText": "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n",
  "LMAfter": "\tfmt.Println(\"Hello\")"
}`,
		},
		{
			name: "JSONL",
			input: `{"FilePath": "main.go", "BaseText": "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n", "LMAfter": "\tfmt.Println(\"Hello\")"}
{"FilePath": "test.go", "BaseText": "package main\n\nfunc test() {\n\treturn\n}\n", "LMAfter": "func Test()"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert([]byte(tt.input), "reviewdog")
			if err != nil {
				if !strings.Contains(err.Error(), "LMAfter") && !strings.Contains(err.Error(), "BaseText") {
					t.Errorf("Convert() unexpected error = %v", err)
				}
				return
			}

			var output interface{}
			if err := json.Unmarshal(result, &output); err != nil {
				t.Errorf("Convert() produced invalid JSON: %v", err)
			}
		})
	}
}
