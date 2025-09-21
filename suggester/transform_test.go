package suggester

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestTransform_Execute(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		core      *CoreResult
		config    *TransformConfig
		wantPaths map[string]interface{} // Expected values at specific paths
	}{
		{
			name: "simple field mapping",
			input: map[string]interface{}{
				"message":    "Fix typo",
				"severity":   "WARNING",
				"sourceName": "test-tool",
				"filePath":   "test.go",
			},
			core: &CoreResult{
				FilePath:    "test.go",
				StartLine:   10,
				StartColumn: 5,
				EndLine:     10,
				EndColumn:   10,
				Before:      "oldText",
				After:       "newText",
			},
			config: &TransformConfig{
				Template: json.RawMessage(`{"result": {}}`),
				Mappings: []FieldMapping{
					{Target: "$.result.message", Source: "input.message"},
					{Target: "$.result.line", Source: "core.StartLine"},
					{Target: "$.result.text", Source: "core.After"},
				},
			},
			wantPaths: map[string]interface{}{
				"result.message": "Fix typo",
				"result.line":    10,
				"result.text":    "newText",
			},
		},
		{
			name: "array and nested mapping",
			input: map[string]interface{}{
				"filePath": "src/main.go",
			},
			core: &CoreResult{
				StartLine:   20,
				StartColumn: 10,
				EndLine:     25,
				EndColumn:   15,
			},
			config: &TransformConfig{
				Template: json.RawMessage(`{"files": [], "meta": {}}`),
				Mappings: []FieldMapping{
					{Target: "$.files[0].path", Source: "input.filePath"},
					{Target: "$.files[0].range.start.line", Source: "core.StartLine"},
					{Target: "$.files[0].range.end.line", Source: "core.EndLine"},
					{Target: "$.meta.version", Source: "literal.1.0.0"},
				},
			},
			wantPaths: map[string]interface{}{
				"files[0].path":            "src/main.go",
				"files[0].range.start.line": 20,
				"files[0].range.end.line":   25,
				"meta.version":              "1.0.0",
			},
		},
		{
			name: "default values",
			input: map[string]interface{}{
				"filePath": "test.js",
			},
			core: &CoreResult{},
			config: &TransformConfig{
				Template: json.RawMessage(`{"tool": {}}`),
				Mappings: []FieldMapping{
					{Target: "$.tool.name", Source: "input.sourceName", Default: "default-tool"},
					{Target: "$.tool.severity", Source: "input.severity", Default: "INFO"},
				},
			},
			wantPaths: map[string]interface{}{
				"tool.name":     "default-tool",
				"tool.severity": "INFO",
			},
		},
		{
			name: "transformation functions",
			input: map[string]interface{}{
				"severity": "ERROR",
			},
			core: &CoreResult{},
			config: &TransformConfig{
				Template: json.RawMessage(`{"sarif": {}, "eslint": {}}`),
				Mappings: []FieldMapping{
					{Target: "$.sarif.level", Source: "input.severity", Transform: "to_sarif_level"},
					{Target: "$.eslint.severity", Source: "input.severity", Transform: "to_eslint_severity"},
				},
			},
			wantPaths: map[string]interface{}{
				"sarif.level":      "error",
				"eslint.severity":  2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transform := &Transform{
				Input:  tt.input,
				Core:   tt.core,
				Config: tt.config,
			}

			got, err := transform.Execute()
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			// Parse result
			var result map[string]interface{}
			if err := json.Unmarshal(got, &result); err != nil {
				t.Fatalf("Failed to parse result JSON: %v", err)
			}

			// Check expected paths
			for path, want := range tt.wantPaths {
				got := getValueByPath(result, path)
				// Special handling for numeric values comparison
				switch wantVal := want.(type) {
				case int:
					switch gotVal := got.(type) {
					case float64:
						if int(gotVal) != wantVal {
							t.Errorf("Path %s: got %v (type %T), want %v (type %T)", path, got, got, want, want)
						}
					case int:
						if gotVal != wantVal {
							t.Errorf("Path %s: got %v, want %v", path, got, want)
						}
					default:
						t.Errorf("Path %s: got %v (type %T), want %v (type %T)", path, got, got, want, want)
					}
				default:
					if got != want {
						t.Errorf("Path %s: got %v (type %T), want %v (type %T)", path, got, got, want, want)
					}
				}
			}
		})
	}
}

func TestParseJSONPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []pathComponent
	}{
		{
			name: "simple field",
			path: "field",
			want: []pathComponent{{Type: "field", Name: "field"}},
		},
		{
			name: "nested fields",
			path: "parent.child",
			want: []pathComponent{
				{Type: "field", Name: "parent"},
				{Type: "field", Name: "child"},
			},
		},
		{
			name: "array index",
			path: "array[0]",
			want: []pathComponent{
				{Type: "field", Name: "array"},
				{Type: "array", Index: 0},
			},
		},
		{
			name: "complex path",
			path: "data.items[2].value",
			want: []pathComponent{
				{Type: "field", Name: "data"},
				{Type: "field", Name: "items"},
				{Type: "array", Index: 2},
				{Type: "field", Name: "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseJSONPath(tt.path)
			if len(got) != len(tt.want) {
				t.Fatalf("parseJSONPath() returned %d components, want %d", len(got), len(tt.want))
			}
			for i, component := range got {
				if component.Type != tt.want[i].Type {
					t.Errorf("Component %d: Type = %s, want %s", i, component.Type, tt.want[i].Type)
				}
				if component.Name != tt.want[i].Name {
					t.Errorf("Component %d: Name = %s, want %s", i, component.Name, tt.want[i].Name)
				}
				if component.Index != tt.want[i].Index {
					t.Errorf("Component %d: Index = %d, want %d", i, component.Index, tt.want[i].Index)
				}
			}
		})
	}
}

// Helper function to get value by path from a map
func getValueByPath(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for i, part := range parts {
		// Handle array notation
		if idx := strings.Index(part, "["); idx >= 0 {
			fieldName := part[:idx]
			endIdx := strings.Index(part, "]")
			indexStr := part[idx+1 : endIdx]
			index := 0
			if num, err := strconv.Atoi(indexStr); err == nil {
				index = num
			}

			// First access the field if present
			if fieldName != "" {
				switch v := current.(type) {
				case map[string]interface{}:
					field, ok := v[fieldName]
					if !ok {
						return nil
					}
					current = field
				default:
					return nil
				}
			}

			// Then access the array index
			switch v := current.(type) {
			case []interface{}:
				if index >= len(v) {
					return nil
				}
				current = v[index]
			default:
				return nil
			}

			// Handle field after array index (e.g., [0].field)
			if endIdx+1 < len(part) && part[endIdx+1] == '.' {
				remainingPath := part[endIdx+2:]
				remainingParts := []string{remainingPath}
				if i+1 < len(parts) {
					remainingParts = append(remainingParts, parts[i+1:]...)
				}
				switch v := current.(type) {
				case map[string]interface{}:
					return getValueByPath(v, strings.Join(remainingParts, "."))
				default:
					return nil
				}
			}
		} else {
			// Simple field access
			switch v := current.(type) {
			case map[string]interface{}:
				next, ok := v[part]
				if !ok {
					return nil
				}
				current = next
			default:
				return nil
			}
		}
	}

	return current
}