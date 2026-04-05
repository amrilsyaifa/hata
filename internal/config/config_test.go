package config

import (
	"os"
	"path/filepath"
	"testing"
)

func sampleConfig() *Config {
	return &Config{
		ProjectID: "test-project",
		Sheet: SheetConf{
			ID:   "sheet-id-123",
			Name: "Translations",
		},
		Auth: AuthConf{
			Type:            "oauth",
			CredentialsPath: ".i18n/credentials.json",
			TokenPath:       ".i18n/token.json",
		},
		Languages: []string{"en-US", "id-ID"},
		Aliases:   map[string]string{"en-US": "en", "id-ID": "id"},
		Paths: PathConf{
			Base:   "./base.json",
			Output: "./locales",
		},
		Options: Options{
			NestedJSON: true,
			SortKeys:   true,
			KeepUnused: false,
		},
	}
}

func TestSaveAndLoad_roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "i18n.config.yml")

	original := sampleConfig()
	if err := Save(path, original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.ProjectID != original.ProjectID {
		t.Errorf("ProjectID: got %q, want %q", loaded.ProjectID, original.ProjectID)
	}
	if loaded.Sheet.ID != original.Sheet.ID {
		t.Errorf("Sheet.ID: got %q, want %q", loaded.Sheet.ID, original.Sheet.ID)
	}
	if loaded.Auth.Type != original.Auth.Type {
		t.Errorf("Auth.Type: got %q, want %q", loaded.Auth.Type, original.Auth.Type)
	}
	if len(loaded.Languages) != len(original.Languages) {
		t.Errorf("Languages length: got %d, want %d", len(loaded.Languages), len(original.Languages))
	}
	if loaded.Options.NestedJSON != original.Options.NestedJSON {
		t.Errorf("Options.NestedJSON: got %v, want %v", loaded.Options.NestedJSON, original.Options.NestedJSON)
	}
	if loaded.Aliases["en-US"] != "en" {
		t.Errorf("Aliases[en-US]: got %q, want %q", loaded.Aliases["en-US"], "en")
	}
}

func TestLoad_missingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_invalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yml")
	if err := os.WriteFile(path, []byte(":	bad:	yaml:"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestSave_createsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "i18n.config.yml")
	if err := Save(path, sampleConfig()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
