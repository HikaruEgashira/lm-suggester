package suggester

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Transform represents a transformation context containing input data,
// core calculation results, and configuration.
type Transform struct {
	Input  map[string]interface{} // Original input as map (for pass-through)
	Core   *CoreResult            // Core calculation results
	Config *TransformConfig       // Transformation configuration
}

// TransformConfig defines how to transform input to output JSON.
type TransformConfig struct {
	Template json.RawMessage `json:"template"` // Output JSON structure template
	Mappings []FieldMapping  `json:"mappings"` // Field mapping definitions
}

// FieldMapping defines how to map a value to output JSON.
type FieldMapping struct {
	Target    string `json:"target"`              // Output JSONPath (e.g., "$.diagnostics[0].location.range.start.line")
	Source    string `json:"source"`              // Data source (e.g., "core.startLine", "input.message")
	Default   string `json:"default,omitempty"`   // Default value if source is empty
	Transform string `json:"transform,omitempty"` // Optional transformation function name
}

// Execute performs the transformation based on configuration.
func (t *Transform) Execute() ([]byte, error) {
	// Parse template to create base output structure
	var output interface{}
	if err := json.Unmarshal(t.Config.Template, &output); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Apply mappings
	for _, mapping := range t.Config.Mappings {
		value, err := t.resolveSource(mapping.Source)
		if err != nil {
			// Use default if source resolution fails
			if mapping.Default != "" {
				value = mapping.Default
			} else {
				continue // Skip if no default and resolution failed
			}
		}

		// Apply transformation if specified
		if mapping.Transform != "" {
			value = t.applyTransform(mapping.Transform, value)
		}

		// Set value in output using JSONPath
		if err := setJSONPath(&output, mapping.Target, value); err != nil {
			return nil, fmt.Errorf("failed to set path %s: %w", mapping.Target, err)
		}
	}

	return json.Marshal(output)
}

// resolveSource resolves a source path to its value.
func (t *Transform) resolveSource(source string) (interface{}, error) {
	parts := strings.SplitN(source, ".", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid source format: %s", source)
	}

	switch parts[0] {
	case "input":
		return getMapValue(t.Input, parts[1])
	case "core":
		return getCoreValue(t.Core, parts[1])
	case "literal":
		return parts[1], nil
	case "computed":
		return t.computeValue(parts[1])
	default:
		return nil, fmt.Errorf("unknown source prefix: %s", parts[0])
	}
}

// getCoreValue gets a value from CoreResult by field name.
func getCoreValue(core *CoreResult, field string) (interface{}, error) {
	if core == nil {
		return nil, fmt.Errorf("core result is nil")
	}

	v := reflect.ValueOf(*core)
	f := v.FieldByName(field)
	if !f.IsValid() {
		return nil, fmt.Errorf("field %s not found in core result", field)
	}
	return f.Interface(), nil
}

// getMapValue gets a nested value from a map using dot notation.
func getMapValue(m map[string]interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := m

	for i, part := range parts {
		if i == len(parts)-1 {
			val, ok := current[part]
			if !ok {
				return nil, fmt.Errorf("key %s not found", part)
			}
			return val, nil
		}

		next, ok := current[part]
		if !ok {
			return nil, fmt.Errorf("key %s not found", part)
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("path %s is not a map", part)
		}
		current = nextMap
	}

	return nil, fmt.Errorf("invalid path: %s", path)
}

// computeValue computes dynamic values.
func (t *Transform) computeValue(name string) (interface{}, error) {
	switch name {
	case "suggestion_message":
		// Format as reviewdog suggestion
		message, _ := getMapValue(t.Input, "Message")
		if message == nil || message == "" {
			message = "Replace code with suggestion"
		}
		return fmt.Sprintf("%s\n```suggestion\n%s\n```", message, t.Core.After), nil
	default:
		return nil, fmt.Errorf("unknown computed value: %s", name)
	}
}

// applyTransform applies a transformation function to a value.
func (t *Transform) applyTransform(name string, value interface{}) interface{} {
	switch name {
	case "to_sarif_level":
		// Convert severity to SARIF level
		switch strings.ToUpper(fmt.Sprint(value)) {
		case "ERROR":
			return "error"
		case "WARNING":
			return "warning"
		case "INFO":
			return "note"
		default:
			return "warning"
		}
	case "to_eslint_severity":
		// Convert to ESLint severity (1=warning, 2=error)
		switch strings.ToUpper(fmt.Sprint(value)) {
		case "ERROR":
			return 2
		case "WARNING":
			return 1
		default:
			return 1
		}
	default:
		return value
	}
}