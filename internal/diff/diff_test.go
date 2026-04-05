package diff

import (
	"testing"

	"github.com/amrilsyaifa/hata/internal/sheet"
)

func rows(keys ...string) []sheet.Row {
	r := make([]sheet.Row, len(keys))
	for i, k := range keys {
		r[i] = sheet.Row{Key: k}
	}
	return r
}

func TestCompare_allInSync(t *testing.T) {
	base := map[string]string{"a": "1", "b": "2"}
	result := Compare(base, rows("a", "b"))
	if len(result.MissingInSheet) != 0 {
		t.Errorf("unexpected missing: %v", result.MissingInSheet)
	}
	if len(result.UnusedInBase) != 0 {
		t.Errorf("unexpected unused: %v", result.UnusedInBase)
	}
}

func TestCompare_missingInSheet(t *testing.T) {
	base := map[string]string{"a": "1", "b": "2", "c": "3"}
	result := Compare(base, rows("a"))
	if len(result.MissingInSheet) != 2 {
		t.Errorf("expected 2 missing, got %v", result.MissingInSheet)
	}
	// Results are sorted
	if result.MissingInSheet[0] != "b" || result.MissingInSheet[1] != "c" {
		t.Errorf("unexpected order: %v", result.MissingInSheet)
	}
}

func TestCompare_unusedInBase(t *testing.T) {
	base := map[string]string{"a": "1"}
	result := Compare(base, rows("a", "b", "c"))
	if len(result.UnusedInBase) != 2 {
		t.Errorf("expected 2 unused, got %v", result.UnusedInBase)
	}
	if result.UnusedInBase[0] != "b" || result.UnusedInBase[1] != "c" {
		t.Errorf("unexpected order: %v", result.UnusedInBase)
	}
}

func TestCompare_empty(t *testing.T) {
	result := Compare(map[string]string{}, rows())
	if len(result.MissingInSheet) != 0 || len(result.UnusedInBase) != 0 {
		t.Errorf("expected empty result, got %+v", result)
	}
}

func TestCompare_sorted(t *testing.T) {
	base := map[string]string{"z": "1", "a": "2", "m": "3"}
	result := Compare(base, rows())
	prev := ""
	for _, k := range result.MissingInSheet {
		if k < prev {
			t.Errorf("results not sorted: %v", result.MissingInSheet)
		}
		prev = k
	}
}
