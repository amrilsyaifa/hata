package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func FlatToNested(flat map[string]string, sortKeys bool) map[string]interface{} {
	result := make(map[string]interface{})

	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}
	if sortKeys {
		sort.Strings(keys)
	}

	for _, key := range keys {
		value := flat[key]
		parts := strings.Split(key, ".")
		current := result
		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = value
			} else {
				if existing, ok := current[part]; ok {
					if nestedMap, ok := existing.(map[string]interface{}); ok {
						current = nestedMap
					} else {
						newMap := make(map[string]interface{})
						current[part] = newMap
						current = newMap
					}
				} else {
					newMap := make(map[string]interface{})
					current[part] = newMap
					current = newMap
				}
			}
		}
	}
	return result
}

func GenerateLocaleFiles(outputDir string, langData map[string]map[string]string, sortKeys, nestedJSON bool) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for lang, flat := range langData {
		var toMarshal interface{}
		if nestedJSON {
			toMarshal = FlatToNested(flat, sortKeys)
		} else {
			if sortKeys {
				keys := make([]string, 0, len(flat))
				for k := range flat {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				ordered := make(map[string]string, len(flat))
				for _, k := range keys {
					ordered[k] = flat[k]
				}
				toMarshal = ordered
			} else {
				toMarshal = flat
			}
		}
		data, err := json.MarshalIndent(toMarshal, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal %s: %w", lang, err)
		}
		path := filepath.Join(outputDir, lang+".json")
		if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
		fmt.Printf("  Generated: %s\n", path)
	}
	return nil
}
