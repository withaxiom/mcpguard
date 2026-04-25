package fingerprint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/withaxiom/mcpguard/internal/config"
)

func TestFileHash(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.json")
	if err := os.WriteFile(f, []byte(`{"test": true}`), 0644); err != nil {
		t.Fatal(err)
	}

	hash, err := FileHash(f)
	if err != nil {
		t.Fatalf("FileHash failed: %v", err)
	}

	if len(hash) != 64 {
		t.Errorf("expected 64-char hex hash, got %d chars", len(hash))
	}

	// Same content should produce same hash
	hash2, _ := FileHash(f)
	if hash != hash2 {
		t.Error("same file should produce same hash")
	}
}

func TestConfigHash_Deterministic(t *testing.T) {
	cfg := &config.MCPConfig{
		Source: "test",
		Path:   "/test",
		Servers: map[string]config.MCPServer{
			"b-server": {Name: "b-server", Command: "node", Args: []string{"b.js"}},
			"a-server": {Name: "a-server", Command: "node", Args: []string{"a.js"}},
		},
	}

	hash1, err := ConfigHash(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Create same config with servers in different order (maps are unordered)
	hash2, err := ConfigHash(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if hash1 != hash2 {
		t.Error("ConfigHash should be deterministic regardless of map order")
	}
}

func TestRedactValue(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"short", "shor*"},
		{"ab", "****"},
		{"ghp_1234567890abcdef", "ghp_****************"},
		{"", "****"},
	}

	for _, tt := range tests {
		got := RedactValue(tt.input)
		if got != tt.want {
			t.Errorf("RedactValue(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsSensitiveKey(t *testing.T) {
	sensitive := []string{
		"API_KEY", "GITHUB_TOKEN", "DB_PASSWORD", "SECRET_KEY",
		"OPENAI_API_KEY", "auth_token", "private_key", "ACCESS_KEY_ID",
	}
	for _, k := range sensitive {
		if !IsSensitiveKey(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}

	safe := []string{
		"NODE_ENV", "HOME", "PATH", "PORT", "DEBUG", "LOG_LEVEL",
	}
	for _, k := range safe {
		if IsSensitiveKey(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}
