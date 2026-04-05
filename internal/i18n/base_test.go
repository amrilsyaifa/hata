package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadBase_valid(t *testing.T) {
	dir := t.TempDir()
	content := `{"hello":"world","foo":"bar"}`
	path := filepath.Join(dir, "base.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	result, err := ReadBase(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["hello"] != "world" {
		t.Errorf("expected hello=world, got %v", result["hello"])
	}
	if result["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", result["foo"])
	}
}

func TestReadBase_missingFile(t *testing.T) {
	_, err := ReadBase("/nonexistent/path/base.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadBase_invalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`not json`), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := ReadBase(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]string{"banana": "b", "apple": "a", "cherry": "c"}
	keys := SortedKeys(m)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "apple" || keys[1] != "banana" || keys[2] != "cherry" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSortedKeys_empty(t *testing.T) {
	keys := SortedKeys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}
