package suggester

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStandardComputeCore_FilePathValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr string
	}{
		{
			name: "FilePath is required when missing",
			input: map[string]interface{}{
				"BaseText": "test content",
				"LMAfter":  "modified content",
			},
			wantErr: "FilePath/file_path is required",
		},
		{
			name: "FilePath is required when empty",
			input: map[string]interface{}{
				"FilePath": "",
				"BaseText": "test content",
				"LMAfter":  "modified content",
			},
			wantErr: "FilePath/file_path is required",
		},
		{
			name: "file_path (lowercase) is also accepted",
			input: map[string]interface{}{
				"file_path": "",
				"base_text": "test content",
				"lm_after":  "modified content",
			},
			wantErr: "FilePath/file_path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StandardComputeCore(tt.input)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestStandardComputeCore_BaseTextFromFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	testContent := `package main

func main() {
	println("hello")
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		input     map[string]interface{}
		wantErr   bool
		errContains string
	}{
		{
			name: "BaseText loaded from file when not provided",
			input: map[string]interface{}{
				"FilePath": testFile,
				"LMBefore": "println(\"hello\")",
				"LMAfter":  "fmt.Println(\"hello\")",
			},
			wantErr: false,
		},
		{
			name: "BaseText loaded from file_path (lowercase)",
			input: map[string]interface{}{
				"file_path": testFile,
				"lm_before": "println(\"hello\")",
				"lm_after":  "fmt.Println(\"hello\")",
			},
			wantErr: false,
		},
		{
			name: "Error when file doesn't exist",
			input: map[string]interface{}{
				"FilePath": "/nonexistent/file.go",
				"LMAfter":  "test",
			},
			wantErr:     true,
			errContains: "failed to read file",
		},
		{
			name: "BaseText from input takes precedence over file",
			input: map[string]interface{}{
				"FilePath": testFile,
				"BaseText": "custom content",
				"LMAfter":  "modified",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StandardComputeCore(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errContains)
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected non-nil result")
			}

			// Verify FilePath is set correctly
			expectedPath := ""
			if fp, ok := tt.input["FilePath"].(string); ok {
				expectedPath = fp
			} else if fp, ok := tt.input["file_path"].(string); ok {
				expectedPath = fp
			}

			if result.FilePath != expectedPath {
				t.Errorf("expected FilePath %q, got %q", expectedPath, result.FilePath)
			}
		})
	}
}

func TestPassthroughConvert_FilePathRequired(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		format  string
		wantErr string
	}{
		{
			name: "reviewdog format requires FilePath",
			input: `{
				"BaseText": "test content",
				"LMAfter": "modified"
			}`,
			format:  "reviewdog",
			wantErr: "FilePath/file_path is required",
		},
		{
			name: "sarif format requires FilePath",
			input: `{
				"BaseText": "test content",
				"LMAfter": "modified"
			}`,
			format:  "sarif",
			wantErr: "FilePath/file_path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := PassthroughConvert([]byte(tt.input), tt.format)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestPassthroughConvert_WithFileRead(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "example.js")
	testContent := `function hello() {
    console.log("world");
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	input := map[string]interface{}{
		"file_path": testFile,
		"lm_before": "console.log(\"world\")",
		"lm_after":  "console.info(\"world\")",
		"message":   "Use console.info for info messages",
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	result, err := PassthroughConvert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var output map[string]interface{}
	if err := json.Unmarshal(result, &output); err != nil {
		t.Fatal(err)
	}

	// Check that the diagnostic was created with correct path
	diagnostics, ok := output["diagnostics"].([]interface{})
	if !ok || len(diagnostics) == 0 {
		t.Fatal("expected diagnostics in output")
	}

	firstDiag, ok := diagnostics[0].(map[string]interface{})
	if !ok {
		t.Fatal("expected diagnostic to be a map")
	}

	location, ok := firstDiag["location"].(map[string]interface{})
	if !ok {
		t.Fatal("expected location in diagnostic")
	}

	path, ok := location["path"].(string)
	if !ok {
		t.Fatal("expected path in location")
	}

	if path != testFile {
		t.Errorf("expected path %q, got %q", testFile, path)
	}
}