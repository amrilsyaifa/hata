package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestSaveToken_andLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")

	tok := &oauth2.Token{
		AccessToken:  "ya29.test-token",
		RefreshToken: "1//refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	}

	if err := saveToken(path, tok); err != nil {
		t.Fatalf("saveToken failed: %v", err)
	}

	loaded, err := loadToken(path)
	if err != nil {
		t.Fatalf("loadToken failed: %v", err)
	}
	if loaded.AccessToken != tok.AccessToken {
		t.Errorf("AccessToken: got %q, want %q", loaded.AccessToken, tok.AccessToken)
	}
	if loaded.RefreshToken != tok.RefreshToken {
		t.Errorf("RefreshToken: got %q, want %q", loaded.RefreshToken, tok.RefreshToken)
	}
}

func TestLoadToken_missingFile(t *testing.T) {
	_, err := loadToken("/nonexistent/token.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadToken_invalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := loadToken(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSaveToken_createsDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "token.json")
	tok := &oauth2.Token{AccessToken: "test", Expiry: time.Now().Add(time.Hour)}
	if err := saveToken(path, tok); err != nil {
		t.Fatalf("saveToken should create directories: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file should exist: %v", err)
	}
}

func TestSavingTokenSource_persistsRefresh(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")

	newTok := &oauth2.Token{
		AccessToken: "new-access-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	// mockTokenSource always returns newTok
	src := &savingTokenSource{
		src:      &staticTokenSource{tok: newTok},
		path:     path,
		lastSeen: "old-token",
	}

	got, err := src.Token()
	if err != nil {
		t.Fatalf("Token() failed: %v", err)
	}
	if got.AccessToken != "new-access-token" {
		t.Errorf("wrong access token: %v", got.AccessToken)
	}

	// File should have been written with the new token
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("token file not created: %v", err)
	}
	var saved oauth2.Token
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("invalid JSON in token file: %v", err)
	}
	if saved.AccessToken != "new-access-token" {
		t.Errorf("saved token mismatch: %v", saved.AccessToken)
	}
}

func TestSavingTokenSource_noWriteWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")

	tok := &oauth2.Token{AccessToken: "same-token", Expiry: time.Now().Add(time.Hour)}
	src := &savingTokenSource{
		src:      &staticTokenSource{tok: tok},
		path:     path,
		lastSeen: "same-token", // same as what the source returns
	}

	if _, err := src.Token(); err != nil {
		t.Fatalf("Token() failed: %v", err)
	}

	// File should NOT be created since token didn't change
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("token file should not have been written when token unchanged")
	}
}

// staticTokenSource is a minimal oauth2.TokenSource for testing.
type staticTokenSource struct {
	tok *oauth2.Token
}

func (s *staticTokenSource) Token() (*oauth2.Token, error) {
	return s.tok, nil
}
