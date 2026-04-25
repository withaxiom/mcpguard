package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClaudeDesktopParser(t *testing.T) {
	// Create a temp config file
	dir := t.TempDir()
	configFile := filepath.Join(dir, "claude_desktop_config.json")
	content := `{
		"mcpServers": {
			"filesystem": {
				"command": "npx",
				"args": ["-y", "@modelcontextprotocol/server-filesystem", "/home/user"],
				"env": {
					"NODE_ENV": "production"
				}
			},
			"github": {
				"command": "npx",
				"args": ["-y", "@modelcontextprotocol/server-github"],
				"env": {
					"GITHUB_TOKEN": "ghp_xxxxxxxxxxxxxxxxxxxx"
				}
			}
		}
	}`
	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &ClaudeDesktopParser{}
	cfg, err := p.Parse(configFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Source != "claude-desktop" {
		t.Errorf("expected source 'claude-desktop', got '%s'", cfg.Source)
	}
	if len(cfg.Servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(cfg.Servers))
	}

	fs, ok := cfg.Servers["filesystem"]
	if !ok {
		t.Fatal("missing 'filesystem' server")
	}
	if fs.Command != "npx" {
		t.Errorf("expected command 'npx', got '%s'", fs.Command)
	}
	if len(fs.Args) != 3 {
		t.Errorf("expected 3 args, got %d", len(fs.Args))
	}
}

func TestCursorParser(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "mcp.json")
	content := `{
		"mcpServers": {
			"database": {
				"command": "node",
				"args": ["server.js"],
				"url": "http://localhost:3000",
				"disabled": true
			}
		}
	}`
	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &CursorParser{}
	cfg, err := p.Parse(configFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Source != "cursor" {
		t.Errorf("expected source 'cursor', got '%s'", cfg.Source)
	}
	if len(cfg.Servers) != 1 {
		t.Errorf("expected 1 server, got %d", len(cfg.Servers))
	}

	db := cfg.Servers["database"]
	if !db.Disabled {
		t.Error("expected database server to be disabled")
	}
	if db.URL != "http://localhost:3000" {
		t.Errorf("expected URL 'http://localhost:3000', got '%s'", db.URL)
	}
}

func TestStripJSONComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "line comment",
			input: `{"key": "value" // comment}`,
			want:  `{"key": "value" `,
		},
		{
			name:  "block comment",
			input: `{"key": /* comment */ "value"}`,
			want:  `{"key":  "value"}`,
		},
		{
			name:  "comment in string preserved",
			input: `{"key": "value // not a comment"}`,
			want:  `{"key": "value // not a comment"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(stripJSONComments([]byte(tt.input)))
			if got != tt.want {
				t.Errorf("stripJSONComments(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
