package suggester

import (
	"encoding/json"
	"testing"
)

func TestConvertWithCustomConfig(t *testing.T) {
	input := Input{
		FilePath:   "test.go",
		BaseText:   "func main() {\n\tfmt.Println(\"Hello\")\n}\n",
		LMBefore:  "fmt.Println(\"Hello\")",
		LMAfter:   "fmt.Println(\"Hello, World!\")",
		Message:    "Add greeting",
		Severity:   "INFO",
		SourceName: "test-linter",
	}

	// Custom config for a simple format
	config := `{
		"template": {
			"file": null,
			"suggestions": []
		},
		"mappings": [
			{"target": "$.file", "source": "input.FilePath"},
			{"target": "$.suggestions[0].line", "source": "core.StartLine"},
			{"target": "$.suggestions[0].message", "source": "input.Message"},
			{"target": "$.suggestions[0].replacement", "source": "core.After"}
		]
	}`

	result, err := ConvertWithCustomConfig(input, []byte(config))
	if err != nil {
		t.Fatalf("ConvertWithCustomConfig failed: %v", err)
	}

	var output map[string]interface{}
	if err := json.Unmarshal(result, &output); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	// Verify the output structure
	if output["file"] != "test.go" {
		t.Errorf("Expected file to be 'test.go', got %v", output["file"])
	}

	suggestions, ok := output["suggestions"].([]interface{})
	if !ok || len(suggestions) == 0 {
		t.Fatal("Expected suggestions array with at least one element")
	}

	firstSuggestion, ok := suggestions[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first suggestion to be a map")
	}

	if firstSuggestion["message"] != "Add greeting" {
		t.Errorf("Expected message to be 'Add greeting', got %v", firstSuggestion["message"])
	}
}

func TestConvertToFormat_SARIF(t *testing.T) {
	input := Input{
		FilePath:   "src/app.js",
		BaseText:   "const x = 1;\nconst y = 2;\n",
		LMBefore:  "const y = 2",
		LMAfter:   "let y = 2",
		Message:    "Use let instead of const",
		Severity:   "WARNING",
		SourceName: "eslint",
	}

	// Note: This test assumes the sarif.json config exists
	// For a real test, we'd need to ensure the config is available
	// or use ConvertWithCustomConfig with an inline SARIF config

	sarifConfig := `{
		"template": {
			"version": "2.1.0",
			"runs": [{
				"tool": {"driver": {"name": null}},
				"results": []
			}]
		},
		"mappings": [
			{"target": "$.runs[0].tool.driver.name", "source": "input.sourceName"},
			{"target": "$.runs[0].results[0].message.text", "source": "input.message"},
			{"target": "$.runs[0].results[0].level", "source": "input.severity", "transform": "to_sarif_level"}
		]
	}`

	result, err := ConvertWithCustomConfig(input, []byte(sarifConfig))
	if err != nil {
		t.Fatalf("Convert to SARIF failed: %v", err)
	}

	var output map[string]interface{}
	if err := json.Unmarshal(result, &output); err != nil {
		t.Fatalf("Failed to parse SARIF result: %v", err)
	}

	// Verify SARIF structure
	if output["version"] != "2.1.0" {
		t.Errorf("Expected SARIF version 2.1.0, got %v", output["version"])
	}

	runs, ok := output["runs"].([]interface{})
	if !ok || len(runs) == 0 {
		t.Fatal("Expected runs array with at least one element")
	}
}

func TestExtractCore(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		baseText string
		lmBefore string
		lmAfter  string
		wantLine int
		wantCol  int
	}{
		{
			name:     "simple replacement",
			filePath: "test.go",
			baseText: "line1\nline2\nline3\n",
			lmBefore: "line2",
			lmAfter:  "modified line2",
			wantLine: 2,
			wantCol:  1,
		},
		{
			name:     "no lmBefore - full diff",
			filePath: "test.py",
			baseText: "def foo():\n    pass\n",
			lmBefore: "",
			lmAfter:  "def foo():\n    return None\n",
			wantLine: 2,
			wantCol:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := ExtractCore(tt.filePath, tt.baseText, tt.lmBefore, tt.lmAfter)
			if err != nil {
				t.Fatalf("ExtractCore failed: %v", err)
			}

			if core.FilePath != tt.filePath {
				t.Errorf("FilePath = %s, want %s", core.FilePath, tt.filePath)
			}

			if core.StartLine != tt.wantLine {
				t.Errorf("StartLine = %d, want %d", core.StartLine, tt.wantLine)
			}

			if core.StartColumn != tt.wantCol {
				t.Errorf("StartColumn = %d, want %d", core.StartColumn, tt.wantCol)
			}
		})
	}
}