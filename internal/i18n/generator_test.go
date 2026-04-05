package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFlatToNested_simple(t *testing.T) {
	flat := map[string]string{
		"greeting": "Hello",
	}
	result := FlatToNested(flat, false)
	if v, ok := result["greeting"]; !ok || v != "Hello" {
		t.Fatalf("expected greeting=Hello, got %v", result)
	}
}

func TestFlatToNested_nested(t *testing.T) {
	flat := map[string]string{
		"auth.login":  "Login",
		"auth.logout": "Logout",
	}
	result := FlatToNested(flat, false)
	auth, ok := result["auth"].(map[string]interface{})
	if !ok {
		t.Fatal("expected auth to be a map")
	}
	if auth["login"] != "Login" {
		t.Errorf("expected auth.login=Login, got %v", auth["login"])
	}
	if auth["logout"] != "Logout" {
		t.Errorf("expected auth.logout=Logout, got %v", auth["logout"])
	}
}

func TestFlatToNested_deeplyNested(t *testing.T) {
	flat := map[string]string{
		"a.b.c.d": "deep",
	}
	result := FlatToNested(flat, false)
	a := result["a"].(map[string]interface{})
	b := a["b"].(map[string]interface{})
	c := b["c"].(map[string]interface{})
	if c["d"] != "deep" {
		t.Errorf("expected a.b.c.d=deep, got %v", c["d"])
	}
}

func TestFlatToNested_empty(t *testing.T) {
	result := FlatToNested(map[string]string{}, false)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestGenerateLocaleFiles_flat(t *testing.T) {
	dir := t.TempDir()
	langData := map[string]map[string]string{
		"en": {"hello": "Hello", "bye": "Bye"},
	}
	if err := GenerateLocaleFiles(dir, langData, true, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "en.json"))
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	var out map[string]string
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if out["hello"] != "Hello" {
		t.Errorf("expected hello=Hello, got %v", out["hello"])
	}
}

func TestGenerateLocaleFiles_nested(t *testing.T) {
	dir := t.TempDir()
	langData := map[string]map[string]string{
		"en": {"auth.login": "Login"},
	}
	if err := GenerateLocaleFiles(dir, langData, true, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "en.json"))
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	auth, ok := out["auth"].(map[string]interface{})
	if !ok {
		t.Fatal("expected nested auth object")
	}
	if auth["login"] != "Login" {
		t.Errorf("expected auth.login=Login, got %v", auth["login"])
	}
}
