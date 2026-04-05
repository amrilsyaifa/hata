package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// ReadBase reads base.json and returns a flat map of dot-notation keys to values.
// It supports both flat JSON ({"key": "value"}) and nested JSON ({"a": {"b": "value"}}).
func ReadBase(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Try flat first
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}

	flat := make(map[string]string)
	if err := flattenJSON("", raw, flat); err != nil {
		return nil, fmt.Errorf("failed to flatten %s: %w", path, err)
	}
	return flat, nil
}

// flattenJSON recursively flattens a JSON value into dot-notation keys.
func flattenJSON(prefix string, value interface{}, out map[string]string) error {
	switch v := value.(type) {
	case map[string]interface{}:
		for k, child := range v {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			if err := flattenJSON(key, child, out); err != nil {
				return err
			}
		}
	case string:
		out[prefix] = v
	case bool, float64:
		out[prefix] = fmt.Sprintf("%v", v)
	case nil:
		out[prefix] = ""
	default:
		return fmt.Errorf("unsupported value type at key %q: %T", prefix, value)
	}
	return nil
}

func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
