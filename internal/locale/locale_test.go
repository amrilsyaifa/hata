package locale

import (
	"strings"
	"testing"
)

func TestDisplay_format(t *testing.T) {
	l := Locale{Code: "en-US", Lang: "English", Country: "United States"}
	d := l.Display()
	if !strings.Contains(d, "(en-US)") {
		t.Errorf("expected Display to contain (en-US), got %q", d)
	}
	if !strings.Contains(d, "English") {
		t.Errorf("expected Display to contain English, got %q", d)
	}
}

func TestCodeFromDisplay_roundtrip(t *testing.T) {
	l := Locale{Code: "id-ID", Lang: "Indonesian", Country: "Indonesia"}
	d := l.Display()
	code := CodeFromDisplay(d)
	if code != "id-ID" {
		t.Errorf("expected id-ID, got %q", code)
	}
}

func TestCodeFromDisplay_invalid(t *testing.T) {
	// Should return the input unchanged when format doesn't match
	code := CodeFromDisplay("no parens here")
	if code != "no parens here" {
		t.Errorf("expected original string back, got %q", code)
	}
}

func TestAll_notEmpty(t *testing.T) {
	if len(All) == 0 {
		t.Fatal("locale list should not be empty")
	}
}

func TestAll_containsCommonLocales(t *testing.T) {
	codes := make(map[string]bool, len(All))
	for _, l := range All {
		codes[l.Code] = true
	}
	for _, expected := range []string{"en-US", "id-ID", "zh-CN", "ja-JP", "fr-FR"} {
		if !codes[expected] {
			t.Errorf("expected locale %q to be in All list", expected)
		}
	}
}

func TestAll_uniqueCodes(t *testing.T) {
	seen := make(map[string]bool)
	for _, l := range All {
		if seen[l.Code] {
			t.Errorf("duplicate locale code: %q", l.Code)
		}
		seen[l.Code] = true
	}
}
