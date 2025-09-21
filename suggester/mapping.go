package suggester

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// setJSONPath sets a value in a JSON structure using a JSONPath-like notation.
// Supports paths like: $.field, $.array[0], $.nested.field, $.array[0].field
func setJSONPath(data *interface{}, path string, value interface{}) error {
	// Remove $ prefix if present
	path = strings.TrimPrefix(path, "$.")

	parts := parseJSONPath(path)
	return setNestedValue(data, parts, value)
}

// parseJSONPath parses a JSONPath string into path components.
func parseJSONPath(path string) []pathComponent {
	var components []pathComponent
	parts := strings.Split(path, ".")

	for _, part := range parts {
		// Check for array notation
		if idx := strings.Index(part, "["); idx >= 0 {
			if idx > 0 {
				// Field name before array index
				components = append(components, pathComponent{
					Type: "field",
					Name: part[:idx],
				})
			}

			// Extract array index
			endIdx := strings.Index(part, "]")
			if endIdx < 0 {
				continue
			}

			indexStr := part[idx+1 : endIdx]
			index, err := strconv.Atoi(indexStr)
			if err == nil {
				components = append(components, pathComponent{
					Type:  "array",
					Index: index,
				})
			}

			// Check for field after array
			remainder := part[endIdx+1:]
			if strings.HasPrefix(remainder, ".") && len(remainder) > 1 {
				components = append(components, pathComponent{
					Type: "field",
					Name: remainder[1:],
				})
			}
		} else {
			// Simple field name
			components = append(components, pathComponent{
				Type: "field",
				Name: part,
			})
		}
	}

	return components
}

type pathComponent struct {
	Type  string // "field" or "array"
	Name  string // For field type
	Index int    // For array type
}

// setNestedValue sets a value in a nested structure.
func setNestedValue(data *interface{}, path []pathComponent, value interface{}) error {
	if len(path) == 0 {
		*data = value
		return nil
	}

	// Helper function to set value at path recursively
	var setValue func(current interface{}, pathIdx int) (interface{}, error)
	setValue = func(current interface{}, pathIdx int) (interface{}, error) {
		if pathIdx >= len(path) {
			return value, nil
		}

		component := path[pathIdx]
		isLast := pathIdx == len(path)-1

		switch component.Type {
		case "field":
			// Ensure current is a map
			var m map[string]interface{}
			switch v := current.(type) {
			case map[string]interface{}:
				m = v
			case nil:
				m = make(map[string]interface{})
			default:
				return nil, fmt.Errorf("expected map at path component %d, got %T", pathIdx, current)
			}

			if isLast {
				m[component.Name] = value
			} else {
				// Recurse into the field
				next, err := setValue(m[component.Name], pathIdx+1)
				if err != nil {
					return nil, err
				}
				m[component.Name] = next
			}
			return m, nil

		case "array":
			// Ensure current is an array
			var arr []interface{}
			switch v := current.(type) {
			case []interface{}:
				arr = v
			case nil:
				arr = make([]interface{}, 0)
			default:
				return nil, fmt.Errorf("expected array at path component %d, got %T", pathIdx, current)
			}

			// Extend array if necessary
			for len(arr) <= component.Index {
				arr = append(arr, nil)
			}

			if isLast {
				arr[component.Index] = value
			} else {
				// Recurse into the array element
				next, err := setValue(arr[component.Index], pathIdx+1)
				if err != nil {
					return nil, err
				}
				arr[component.Index] = next
			}
			return arr, nil
		}

		return nil, fmt.Errorf("unknown path component type: %s", component.Type)
	}

	result, err := setValue(*data, 0)
	if err != nil {
		return err
	}
	*data = result
	return nil
}

// LoadConfig loads a TransformConfig from JSON data.
func LoadConfig(data []byte) (*TransformConfig, error) {
	var config TransformConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &config, nil
}

// LoadConfigFromFile loads a TransformConfig from a file.
func LoadConfigFromFile(path string) (*TransformConfig, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return LoadConfig(data)
}

// readFile is a helper to read file contents.
func readFile(path string) ([]byte, error) {
	// This would typically use os.ReadFile
	// For now, returning an error as placeholder
	return nil, fmt.Errorf("file reading not implemented in this context")
}